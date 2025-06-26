package models

import (
	"time"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/types"
	pkg "github.com/Shaheen-AlQaraghuli/wallet-go/pkg/wallet"
	"gorm.io/gorm"
)

type Wallets []Wallet
type Wallet struct {
	ID        string
	OwnerID   string
	Currency  string
	Status    string
	Balance   *int `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (w Wallet) ToResponse() pkg.Wallet {
	return pkg.Wallet{
		ID:        w.ID,
		OwnerID:   w.OwnerID,
		Currency:  types.Currency(w.Currency),
		Status:    types.WalletStatus(w.Status),
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

func (w Wallets) ToResponse() []pkg.Wallet {
	res := make([]pkg.Wallet, 0, len(w))

	for _, wallet := range w {
		res = append(res, wallet.ToResponse())
	}

	return res
}

type QueryWallets struct {
	IDs        []string
	OwnerIDs   []string
	Currencies []string

	pagination.Paginator
}

func (q QueryWallets) FromRequest(req pkg.ListWalletsRequest) QueryWallets {
	return QueryWallets{
		IDs:        req.IDs,
		OwnerIDs:   req.OwnerIDs,
		Currencies: req.Currencies.String(),
		Paginator:  req.Paginator,
	}
}

type CreateWalletRequest struct {
	OwnerID  string
	Currency string
	Status   string
}

func (c CreateWalletRequest) FromRequest(req pkg.CreateWalletRequest) CreateWalletRequest {
	return CreateWalletRequest{
		OwnerID:  req.OwnerID,
		Currency: req.Currency.String(),
		Status:   string(types.WalletStatusActive),
	}
}

func (c CreateWalletRequest) ToWallet() Wallet {
	return Wallet{
		OwnerID:  c.OwnerID,
		Currency: c.Currency,
		Status:   c.Status,
	}
}
