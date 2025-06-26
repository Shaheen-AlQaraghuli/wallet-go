package types

type TransactionStatus string
type TransactionStatuses []TransactionStatus

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

func (t TransactionStatus) String() string {
	return string(t)
}

func (t TransactionStatuses) String() []string {
	strs := make([]string, 0, len(t))
	for _, status := range t {
		strs = append(strs, status.String())
	}

	return strs
}

func GetTransactionStatuses() []TransactionStatus {
	return []TransactionStatus{
		TransactionStatusPending,
		TransactionStatusCompleted,
		TransactionStatusFailed,
	}
}
