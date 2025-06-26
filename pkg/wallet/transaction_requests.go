package wallet

import (
	"time"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/types"
)

//nolint:godox
type CreateTransactionRequest struct {
	// Unique identifier for the wallet.
	WalletID string `binding:"required" form:"wallet_id" json:"wallet_id" url:"wallet_id"`
	// Amount to be added or deducted from the wallet.
	Amount int `binding:"required" form:"amount" json:"amount" url:"amount"`
	// Note for the transaction.
	Note *string `binding:"omitempty" form:"note,omitempty" json:"note,omitempty" url:"note,omitempty"`
	// Type of transaction.
	Type types.TransactionType `binding:"required,transactionTypeEnum" form:"type" json:"type" url:"type"`
	// Idempotency key for the transaction.
	IdempotencyKey string `binding:"required" form:"idempotency_key" json:"idempotency_key" url:"idempotency_key"`
}

//nolint:lll
type ListTransactionsRequest struct {
	// IDs of the transactions to filter.
	IDs []string `binding:"omitempty" form:"ids,omitempty" json:"ids,omitempty" url:"ids,omitempty"`
	// Wallet IDs to filter.
	WalletIDs []string `binding:"omitempty" form:"wallet_ids,omitempty" json:"wallet_ids,omitempty" url:"wallet_ids,omitempty"`
	// Statuses of the transactions to filter.
	Statuses types.TransactionStatuses `binding:"omitempty,transactionStatusesEnum" form:"statuses,omitempty" json:"statuses,omitempty" url:"statuses,omitempty"`
	// Types of the transactions to filter.
	Types types.TransactionTypes `binding:"omitempty,transactionTypesEnum" form:"types,omitempty" json:"types,omitempty" url:"types,omitempty"`
	// CreatedAtFrom is the start date for filtering transactions.
	CreatedAtFrom *time.Time `binding:"omitempty" form:"created_at_from,omitempty" json:"created_at_from,omitempty" url:"created_at_from,omitempty"`
	// CreatedAtTo is the end date for filtering transactions.
	CreatedAtTo *time.Time `binding:"omitempty" form:"created_at_to,omitempty" json:"created_at_to,omitempty" url:"created_at_to,omitempty"`

	pagination.Paginator
}

type UpdateTransactionStatusRequest struct {
	Status types.TransactionStatus `binding:"required,transactionStatusEnum" form:"status" json:"status" url:"status"`
}
