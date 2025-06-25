package wallet

import (
	"time"
)

type Wallet struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	Balance   *int      `json:"balance,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletResponse struct {
	Wallet `json:"wallet"`
}

type WalletsResponse struct {
	Wallets  []Wallet `json:"wallets"`
	Metadata Metadata `json:"metadata"`
}
