package transactions

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
	"wallet/internal/app/models"
	"wallet/internal/util/pagination"
)

type transactionDB interface {
	Create(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
	GetByID(ctx context.Context, id string) (models.Transaction, error)
	Update(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
	List(ctx context.Context, query models.QueryTransactions) ([]models.Transaction, *pagination.Pagination, error)
	ListAllTransactions(ctx context.Context, walletID string) (models.Transactions, error)
}

type walletService interface {
	GetByID(ctx context.Context, id string) (models.Wallet, error)
}

type cacheClient interface {
	GetBalance(ctx context.Context, walletID string) (*int, error)
	SetBalance(ctx context.Context, walletID string, balance int) error
	Mutex(ctx context.Context, key string) (func(), error)
	GetIdempotentTransaction(ctx context.Context, idempotencyKey string) (*models.Transaction, error)
	SetIdempotentTransaction(ctx context.Context, idempotencyKey string, transaction models.Transaction) error
}

type Service struct {
	walletService walletService
	db            transactionDB
	cache         cacheClient
}

func NewService(walletService walletService, db transactionDB, cache cacheClient) *Service {
	return &Service{
		walletService: walletService,
		db:            db,
		cache:         cache,
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

	newStatus := fsm.NewFSM(transaction.Status, models.TransactionStates, nil)
	if newStatus.Cannot(status) {
		return models.Transaction{}, fmt.Errorf("invalid status transition from %s to %s", transaction.Status, status)
	}

	transaction.Status = status

	return s.db.Update(ctx, transaction)
}

func (s *Service) ListTransactions(ctx context.Context, query models.QueryTransactions) ([]models.Transaction, *pagination.Pagination, error) {
	return s.db.List(ctx, query)
}
