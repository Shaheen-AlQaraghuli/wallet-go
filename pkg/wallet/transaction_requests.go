package wallet

import (
	"wallet/pkg/types"
)

type CreateTransactionRequest struct {
	// Unique identifier for the wallet.
	WalletID string `binding:"required" form:"wallet_id" json:"wallet_id" url:"wallet_id"`
	// Amount to be added or deducted from the wallet.
	Amount int `binding:"required" form:"amount" json:"amount" url:"amount"`
	// Note for the transaction.
	Note *string `binding:"omitempty" form:"note,omitempty" json:"note,omitempty" url:"note,omitempty"`
	// Type of transaction credit or debit.
	Type string `binding:"required" form:"type" json:"type" url:"type"`
	// Status of the transaction.
	Status string `binding:"required" form:"status" json:"status" url:"status"`
	// Idempotency key for the transaction.
	IdempotencyKey string `binding:"required" form:"idempotency_key" json:"idempotency_key" url:"idempotency_key"`
}

type ListTransactionsRequest struct {
	// IDs of the transactions to filter.
	IDs []string `binding:"omitempty" form:"ids,omitempty" json:"ids,omitempty" url:"ids,omitempty"`
	// Wallet IDs to filter.
	WalletIDs []string `binding:"omitempty" form:"wallet_ids,omitempty" json:"wallet_ids,omitempty" url:"wallet_ids,omitempty"`
	// Statuses of the transactions to filter.
	Statuses []types.TransactionStatus `binding:"omitempty,transactionStatusEnum" form:"statuses,omitempty" json:"statuses,omitempty" url:"statuses,omitempty"`
	// Types of the transactions to filter.
	Types []types.TransactionType `binding:"omitempty,transactionTypeEnum" form:"types,omitempty" json:"types,omitempty" url:"types,omitempty"`
	// CreatedAtFrom is the start date for filtering transactions.
	CreatedAtFrom string `binding:"omitempty" form:"created_at_from,omitempty" json:"created_at_from,omitempty" url:"created_at_from,omitempty"`
	// CreatedAtTo is the end date for filtering transactions.
	CreatedAtTo string `binding:"omitempty" form:"created_at_to,omitempty" json:"created_at_to,omitempty" url:"created_at_to,omitempty"`
}

type UpdateTransactionStatusRequest struct {
	Status types.TransactionStatus `binding:"required,transactionStatusEnum" form:"status" json:"status" url:"status"`
}
