package wallet

import "wallet/pkg/types"

type Transaction struct {
	ID        string                  `json:"id"`
	WalletID  string                  `json:"wallet_id"`
	Amount    int                     `json:"amount"`
	Note      *string                 `json:"note,omitempty"`
	Type      types.TransactionType   `json:"type"`
	Status    types.TransactionStatus `json:"status"`
	CreatedAt string                  `json:"created_at"`
	UpdatedAt string                  `json:"updated_at"`
}

type TransactionResponse struct {
	Transaction `json:"transaction"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Metadata     Metadata      `json:"metadata"`
}
