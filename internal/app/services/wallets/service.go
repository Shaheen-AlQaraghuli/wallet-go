package wallets

import (
	"context"
	"errors"
	"time"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/ulid"
)

type walletDB interface {
	Create(ctx context.Context, wallet models.Wallet) (models.Wallet, error)
	GetByID(ctx context.Context, id string) (models.Wallet, error)
	Update(ctx context.Context, wallet models.Wallet) (models.Wallet, error)
	List(ctx context.Context, query models.QueryWallets) ([]models.Wallet, *pagination.Pagination, error)
}

type transactionService interface {
	RunningBalance(ctx context.Context, walletID string) (int, error)
}

type cache interface {
	GetBalance(ctx context.Context, walletID string) (*int, error)
}

type Service struct {
	transactionService transactionService
	db                 walletDB
	cache              cache
	now                func() time.Time
}

func NewService(transactionService transactionService, db walletDB, cache cache, now func() time.Time) *Service {
	return &Service{
		transactionService: transactionService,
		db:                 db,
		cache:              cache,
		now:                now,
	}
}

func (s *Service) GetWalletByID(ctx context.Context, id string) (models.Wallet, error) {
	return s.db.GetByID(ctx, id)
}

func (s *Service) UpdateWalletStatus(ctx context.Context, id, status string) (models.Wallet, error) {
	wallet, err := s.db.GetByID(ctx, id)
	if err != nil {
		return models.Wallet{}, err
	}

	if wallet.Status == status {
		return wallet, nil
	}

	wallet.Status = status

	return s.db.Update(ctx, wallet)
}

func (s *Service) ListWallets(ctx context.Context, query models.QueryWallets) (
	models.Wallets,
	*pagination.Pagination, error,
) {
	return s.db.List(ctx, query)
}

func (s *Service) CreateWallet(ctx context.Context, req models.CreateWalletRequest) (models.Wallet, error) {
	list, _, err := s.db.List(ctx, models.QueryWallets{
		OwnerIDs:   []string{req.OwnerID},
		Currencies: []string{req.Currency},
	})
	if err != nil {
		return models.Wallet{}, err
	}

	if len(list) > 0 {
		return models.Wallet{}, errors.New("wallet already exists for this owner and currency")
	}

	newID := ulid.GenerateID(s.now())
	wallet := req.ToWallet()
	wallet.ID = newID

	return s.db.Create(ctx, wallet)
}

func (s *Service) GetWalletWithBalance(ctx context.Context, id string) (models.Wallet, error) {
	wallet, err := s.db.GetByID(ctx, id)
	if err != nil {
		return models.Wallet{}, err
	}

	balance, err := s.transactionService.RunningBalance(ctx, id)
	if err != nil {
		return models.Wallet{}, err
	}

	wallet.Balance = &balance

	return wallet, nil
}
