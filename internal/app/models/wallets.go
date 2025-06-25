package models

import (
	"time"

	"gorm.io/gorm"
	"wallet/internal/util/pagination"
)

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

type QueryWallets struct {
	IDs        []string
	OwnerIDs   []string
	Currencies []string

	*pagination.Paginator
}
