package transactions

import (
	"context"
	"fmt"
	"log"
	"time"

	"wallet/internal/app/models"
	"wallet/internal/util/dblib"
	"wallet/internal/util/pagination"
	"wallet/pkg/types"

	"github.com/looplab/fsm"
	"go.uber.org/zap"
)

type transactionRepo interface {
	dblib.TxManager

	Create(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
	GetByID(ctx context.Context, id string) (models.Transaction, error)
	Update(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
	List(ctx context.Context, query models.QueryTransactions) ([]models.Transaction, *pagination.Pagination, error)
	ListAllTransactions(ctx context.Context, walletID string) (models.Transactions, error)
}

type walletRepo interface {
	GetByID(ctx context.Context, id string) (models.Wallet, error)
}

type cacheClient interface {
	GetBalance(ctx context.Context, walletID string) (*int, error)
	SetBalance(ctx context.Context, walletID string, balance int) error
	Mutex(ctx context.Context, key string) (func(context.Context) (bool, error), error)
	GetIdempotentTransaction(ctx context.Context, idempotencyKey string) (*models.Transaction, error)
	SetIdempotentTransaction(ctx context.Context, idempotencyKey string, transaction models.Transaction) error
}

type Service struct {
	walletRepo walletRepo
	db         transactionRepo
	cache      cacheClient
	now       func() time.Time
}

func NewService(walletRepo walletRepo, db transactionRepo, cache cacheClient, now func() time.Time) *Service {
	return &Service{
		walletRepo: walletRepo,
		db:        db,
		cache:     cache,
		now:      now,
	}
}

func (s *Service) GetTransactionByID(ctx context.Context, id string) (models.Transaction, error) {
	return s.db.GetByID(ctx, id)
}

func (s *Service) UpdateTransactionStatus(ctx context.Context, id string, status string) (models.Transaction, error) {
	transaction, err := s.db.GetByID(ctx, id)
	if err != nil {
		return models.Transaction{}, err
	}

	if transaction.Status == status {
		return transaction, nil
	}

	unlock, err := s.cache.Mutex(ctx, transaction.WalletID)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed to lock wallet: %w", err)
	}
	defer unlock(ctx)

	newStatus := fsm.NewFSM(transaction.Status, models.TransactionStates, nil)
	if newStatus.Cannot(status) {
		return models.Transaction{}, fmt.Errorf("invalid status transition from %s to %s", transaction.Status, status)
	}

	return s.updateStatus(ctx, transaction, status)
}

func (s *Service) ListTransactions(ctx context.Context, query models.QueryTransactions) (models.Transactions, *pagination.Pagination, error) {
	return s.db.List(ctx, query)
}

func (s *Service) updateStatus(ctx context.Context, transaction models.Transaction, status string) (models.Transaction, error) {
	var (
		updatedTransaction models.Transaction
		err                error
	)

	if err := s.db.Tx(ctx, func(ctx context.Context) error {
		prevStatus := transaction.Status
		transaction.Status = status

		updatedTransaction, err = s.db.Update(ctx, transaction)
		if err != nil {
			return err
		}

		return s.refreshCacheAfterTransactionUpdate(ctx, updatedTransaction, prevStatus)
	}); err != nil {
  			return models.Transaction{}, err
 }

 	return updatedTransaction, nil
}

func (s *Service) refreshCacheAfterTransactionUpdate(ctx context.Context, transaction models.Transaction, previousStatus string) error {
	if !shouldUpdateBalanceCache(transaction, previousStatus) {
  		return nil
 	}

	balance, err := s.cache.GetBalance(ctx, transaction.WalletID)
	if err != nil {
		log.Println("error getting balance from cache:",
			zap.Error(err),
			zap.String("walletID", transaction.WalletID),
			zap.Any("transaction", transaction))
		return err
	}

	//no need to refresh if balance is not in cache
	if balance == nil {
		return nil
	}

	switch transaction.Type {
	case string(types.TransactionTypeCredit):
		*balance -= transaction.Amount
	case string(types.TransactionTypeDebit):
		*balance += transaction.Amount
	}

	if err := s.cache.SetBalance(ctx, transaction.WalletID, *balance); err != nil {
		log.Println("error setting balance in cache:",
			zap.Error(err),
			zap.String("walletID", transaction.WalletID),
			zap.Any("transaction", transaction))
		return err
	}

	return nil
}

func shouldUpdateBalanceCache(transaction models.Transaction, previousStatus string) bool {
	return transaction.Status != previousStatus && transaction.Status == string(types.TransactionStatusFailed)
}
