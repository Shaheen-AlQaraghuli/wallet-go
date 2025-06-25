package transactions

import (
	"context"
	"log"

	"go.uber.org/zap"
)

func (s *Service) RunningBalance(ctx context.Context, walletID string) (int, error) {
	balance, err := s.cache.GetBalance(ctx, walletID)
	if err != nil {
		return 0, err
	}

	if balance != nil {
		return *balance, nil
	}

	unlock, err := s.cache.Mutex(ctx, walletID)
	if err != nil {
		log.Println("error locking wallet:", zap.Error(err), zap.String("walletID", walletID))

		return 0, err
	}

	defer func() {
		unlock()
	}()

	transactions, err := s.db.ListAllTransactions(ctx, walletID)
	if err != nil {
		log.Println("error listing all transactions:", zap.Error(err), zap.String("walletID", walletID))

		return 0, err
	}

	err = s.cache.SetBalance(ctx, walletID, transactions.Balance())
	if err != nil {
		log.Println("error setting balance in cache:", zap.Error(err), zap.String("walletID", walletID))
	}

	return transactions.Balance(), nil
}
