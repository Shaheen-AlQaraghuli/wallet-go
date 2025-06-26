package transactions

import (
	"context"
	"fmt"
	"log"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/types"
	"github.com/looplab/fsm"
	"go.uber.org/zap"
)

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

func (s *Service) updateStatus(
	ctx context.Context,
	transaction models.Transaction,
	status string,
) (models.Transaction, error) {
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

func (s *Service) refreshCacheAfterTransactionUpdate(
	ctx context.Context,
	transaction models.Transaction,
	previousStatus string,
) error {
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

	//no need to refresh if balance is not in cache.
	if balance == nil {
		return nil
	}

	newBalance := computeNewBalanceAfterTransactionStatusUpdate(
		transaction.Status,
		transaction,
		*balance,
	)

	if err := s.cache.SetBalance(ctx, transaction.WalletID, newBalance); err != nil {
		log.Println("error setting balance in cache:",
			zap.Error(err),
			zap.String("walletID", transaction.WalletID),
			zap.Any("transaction", transaction))

		return err
	}

	return nil
}

func shouldUpdateBalanceCache(transaction models.Transaction, previousStatus string) bool {
	//update if current is failed or pending to completed.
	return transaction.Status != previousStatus &&
		(transaction.Status == string(types.TransactionStatusFailed) ||
			previousStatus == string(types.TransactionStatusPending) &&
				transaction.Status == string(types.TransactionStatusCompleted))
}

func computeNewBalanceAfterTransactionStatusUpdate(
	status string,
	transaction models.Transaction,
	currentBalance int,
) int {
	if status == string(types.TransactionStatusFailed) {
		if transaction.Type == string(types.TransactionTypeDebit) {
			return currentBalance + transaction.Amount
		}

		return currentBalance - transaction.Amount
	}

	// TransactionStatusComplted.

	if transaction.Type == string(types.TransactionTypeCredit) {
		return currentBalance + transaction.Amount
	}

	//No need to change for pending debit as it is already handled in create transaction.
	return currentBalance
}
