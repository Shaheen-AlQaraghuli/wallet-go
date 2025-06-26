package wallet

import (
	"time"

	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/types"
)

type Wallet struct {
	ID        string             `json:"id"`
	OwnerID   string             `json:"owner_id"`
	Currency  types.Currency     `json:"currency"`
	Status    types.WalletStatus `json:"status"`
	Balance   *int               `json:"balance,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type WalletResponse struct {
	Wallet `json:"wallet"`
}

type WalletsResponse struct {
	Wallets  []Wallet `json:"wallets"`
	Metadata Metadata `json:"metadata"`
}
