package transactions

import (
	"context"
	"errors"
	"fmt"
	"log"

	"wallet/internal/app/models"
	"wallet/internal/util/ulid"
	"wallet/pkg/types"

	"go.uber.org/zap"
)

func (s *Service) CreateTransaction(ctx context.Context, req models.CreateTransactionRequest) (models.Transaction, error) {
	// Lock on the idempotency key to prevent race conditions.
	idempotencyUnlock, err := s.cache.Mutex(ctx, fmt.Sprintf("idempotency:%s", req.IdempotencyKey))
	if err != nil {
		log.Println("error locking idempotency key:", zap.Error(err), zap.String("idempotencyKey", req.IdempotencyKey))

		return models.Transaction{}, err
	}

	defer func() {
		idempotencyUnlock(ctx)
	}()

	existingTransaction, err := s.cache.GetIdempotentTransaction(ctx, req.IdempotencyKey)
	if err != nil {
		log.Println("error checking idempotency key:", zap.Error(err), zap.String("idempotencyKey", req.IdempotencyKey))

		return models.Transaction{}, err
	}

	if existingTransaction != nil {
		log.Println("returning cached transaction for idempotency key:", zap.String("idempotencyKey", req.IdempotencyKey), zap.String("transactionID", existingTransaction.ID))

		return *existingTransaction, nil
	}

	return s.create(ctx, req)
}

func (s *Service) create(ctx context.Context, req models.CreateTransactionRequest) (models.Transaction, error) {
	wallet, err := s.walletRepo.GetByID(ctx, req.WalletID)
	if err != nil {
		log.Println("error getting wallet by ID:", zap.Error(err), zap.String("walletID", req.WalletID))

		return models.Transaction{}, err
	}

	// lock the wallet to prevent race conditions.
	unlock, err := s.cache.Mutex(ctx, wallet.ID)
	if err != nil {
		log.Println("error locking wallet:", zap.Error(err))

		return models.Transaction{}, err
	}

	defer func() {
		unlock(ctx)
	}()

	//todo move this to another function validate
	ledger, err := s.db.ListAllTransactions(ctx, wallet.ID)
	if err != nil {
		log.Println("error listing all transactions:", zap.Error(err), zap.String("walletID", wallet.ID))

		return models.Transaction{}, err
	}

	if req.Type == string(types.TransactionTypeDebit) && ledger.Balance() < req.Amount {
		log.Println("insufficient funds for transaction:", zap.String("walletID", wallet.ID), zap.Int("transactionAmount", req.Amount), zap.Int("balance", ledger.Balance()))

		return models.Transaction{}, errors.New("insufficient funds")
	}

	transaction := req.ToTransaction()
	transaction.ID = ulid.GenerateID(s.now())

	transaction, err = s.db.Create(ctx, transaction)
	if err != nil {
		log.Println("error creating transaction:", zap.Error(err))

		return models.Transaction{}, err
	}

	if err := s.cache.SetIdempotentTransaction(ctx, req.IdempotencyKey, transaction); err != nil {
		log.Println("error caching transaction for idempotency:", zap.Error(err), zap.String("idempotencyKey", req.IdempotencyKey))
	}

	return s.updateBalanceInCache(ctx, ledger.Balance(), transaction)
}

func (s *Service) updateBalanceInCache(ctx context.Context, currentBalance int, transaction models.Transaction) (models.Transaction, error) {
	if transaction.Status != string(types.TransactionStatusFailed) {
		if transaction.Type == string(types.TransactionTypeDebit) {
			currentBalance -= transaction.Amount
		} else if transaction.Type == string(types.TransactionTypeCredit) {
			currentBalance += transaction.Amount
		}
	}

	if err := s.cache.SetBalance(ctx, transaction.WalletID, currentBalance); err != nil {
		log.Println("error setting balance in cache:", zap.Error(err), zap.String("walletID", transaction.WalletID))
	}

	return transaction, nil
}
