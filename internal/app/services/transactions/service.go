package transactions

import (
	"context"
	"time"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/dblib"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
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
	now        func() time.Time
}

func NewService(walletRepo walletRepo, db transactionRepo, cache cacheClient, now func() time.Time) *Service {
	return &Service{
		walletRepo: walletRepo,
		db:         db,
		cache:      cache,
		now:        now,
	}
}

func (s *Service) GetTransactionByID(ctx context.Context, id string) (models.Transaction, error) {
	return s.db.GetByID(ctx, id)
}

func (s *Service) ListTransactions(
	ctx context.Context,
	query models.QueryTransactions,
) (models.Transactions, *pagination.Pagination, error) {
	return s.db.List(ctx, query)
}
