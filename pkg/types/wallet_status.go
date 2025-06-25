package types

type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "active"
	WalletStatusInactive WalletStatus = "inactive"
	WalletStatusFrozen   WalletStatus = "frozen"
)

func GetWalletStatuses() []WalletStatus {
	return []WalletStatus{
		WalletStatusActive,
		WalletStatusInactive,
		WalletStatusFrozen,
	}
}
