package types

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

func GetTransactionStatuses() []TransactionStatus {
	return []TransactionStatus{
		TransactionStatusPending,
		TransactionStatusCompleted,
		TransactionStatusFailed,
	}
}
