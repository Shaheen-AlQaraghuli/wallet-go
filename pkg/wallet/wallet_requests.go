package wallet

import (
	"wallet/internal/util/pagination"
	types "wallet/pkg/types"
)

type CreateWalletRequest struct {
	// Unique identifier for the wallet owner.
	OwnerID string `binding:"required" form:"owner_id" json:"owner_id" url:"owner_id"`
	// Currency of the wallet.
	Currency types.Currency `binding:"required,currencyEnum" form:"currency" json:"currency" url:"currency"`
}

//nolint:lll
type ListWalletsRequest struct {
	// IDs of the wallets to filter.
	IDs []string `binding:"omitempty" form:"ids,omitempty" json:"ids,omitempty" url:"ids,omitempty"`
	// Owner IDs to filter.
	OwnerIDs []string `binding:"omitempty" form:"owner_ids,omitempty" json:"owner_ids,omitempty" url:"owner_ids,omitempty"`
	// Currencies to filter.
	Currencies types.Currencies `binding:"omitempty,walletStatusesEnum" form:"currencies,omitempty" json:"currencies,omitempty" url:"currencies,omitempty"`

	pagination.Paginator
}

type UpdateWalletStatusRequest struct {
	Status types.WalletStatus `binding:"required,walletStatusEnum" form:"status" json:"status" url:"status"`
}
