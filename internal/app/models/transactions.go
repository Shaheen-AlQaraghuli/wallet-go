package models

import (
	"time"

	"github.com/looplab/fsm"
	"wallet/internal/util/pagination"
	"wallet/pkg/types"
)

type Transaction struct {
	ID            string
	WalletID      string
	Amount        int
	Note          *string
	Type          string
	Status        string
	IdempotencyKey string `gorm:"-"`
	CreatedAt     time.Time 
	UpdatedAt     time.Time 
}

type Transactions []Transaction

type QueryTransactions struct {
	IDs           []string
	WalletIDs     []string
	Statuses      []string
	Types         []string   
	CreatedAtFrom *time.Time 
	CreatedAtTo   *time.Time 

	*pagination.Paginator
}

func (t Transactions) Balance() int {
	balance := 0
	for _, transaction := range t {
		if transaction.Status == string(types.TransactionStatusFailed) {
			continue
		}

		if transaction.Type == string(types.TransactionTypeCredit) {
			balance += int(transaction.Amount)
		} else if transaction.Type == string(types.TransactionTypeDebit) {
			balance -= int(transaction.Amount)
		}
	}

	return balance
}

var (
	TransactionStates = fsm.Events{
		{
			Name: string(types.TransactionStatusPending),
			Src:  []string{""},
			Dst:  string(types.TransactionStatusPending),
		},
		{
			Name: string(types.TransactionStatusCompleted),
   			Src:  []string{string(types.TransactionStatusPending)},
   			Dst:  string(types.TransactionStatusCompleted),
		},
		{
			Name: string(types.TransactionStatusFailed),
			Src:  []string{string(types.TransactionStatusPending)},
			Dst:  string(types.TransactionStatusFailed),
		},
	}
)
