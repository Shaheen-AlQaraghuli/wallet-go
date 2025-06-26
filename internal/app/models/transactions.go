package models

import (
	"time"

	"github.com/looplab/fsm"
	"wallet/internal/util/pagination"
	"wallet/pkg/types"
	pkg "wallet/pkg/wallet"
)

type Transaction struct {
	ID        string
	WalletID  string
	Amount    int
	Note      *string
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transactions []Transaction

func (t Transaction) ToResponse() pkg.Transaction {
	return pkg.Transaction{
		ID:        t.ID,
		WalletID:  t.WalletID,
		Amount:    t.Amount,
		Note:      t.Note,
		Type:      types.TransactionType(t.Type),
		Status:    types.TransactionStatus(t.Status),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (t Transactions) ToResponse() []pkg.Transaction {
	res := make([]pkg.Transaction, 0, len(t))
	for _, transaction := range t {
		res = append(res, transaction.ToResponse())
	}

	return res
}

type QueryTransactions struct {
	IDs           []string
	WalletIDs     []string
	Statuses      []string
	Types         []string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time

	pagination.Paginator
}

func (q QueryTransactions) FromRequest(req pkg.ListTransactionsRequest) QueryTransactions {
	return QueryTransactions{
		IDs:           req.IDs,
		WalletIDs:     req.WalletIDs,
		Statuses:      req.Statuses.String(),
		Types:         req.Types.String(),
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
		Paginator:     req.Paginator,
	}
}

func (t Transactions) Balance() int {
	balance := 0

	for _, transaction := range t {
		if transaction.Status == string(types.TransactionStatusFailed) {
			continue
		}

		if transaction.Type == string(types.TransactionTypeCredit) &&
			transaction.Status != string(types.TransactionStatusPending) { //nolint:wsl

			balance += transaction.Amount
		} else if transaction.Type == string(types.TransactionTypeDebit) {
			balance -= transaction.Amount
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

type CreateTransactionRequest struct {
	WalletID       string
	Amount         int
	Note           *string
	Type           string
	IdempotencyKey string
}

func (r CreateTransactionRequest) FromRequest(req pkg.CreateTransactionRequest) CreateTransactionRequest {
	return CreateTransactionRequest{
		WalletID:       req.WalletID,
		Amount:         req.Amount,
		Note:           req.Note,
		Type:           req.Type.String(),
		IdempotencyKey: req.IdempotencyKey,
	}
}

func (r CreateTransactionRequest) ToTransaction() Transaction {
	return Transaction{
		WalletID: r.WalletID,
		Amount:   r.Amount,
		Note:     r.Note,
		Type:     r.Type,
		Status:   string(types.TransactionStatusPending),
	}
}
