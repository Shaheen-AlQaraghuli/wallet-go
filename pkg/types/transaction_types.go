package types

type TransactionType string
type TransactionTypes []TransactionType

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"
)

func (t TransactionType) String() string {
	return string(t)
}

func (t TransactionTypes) String() []string {
	strs := make([]string, 0, len(t))
	for _, status := range t {
		strs = append(strs, status.String())
	}

	return strs
}

func GetTransactionTypes() []TransactionType {
	return []TransactionType{
		TransactionTypeCredit,
		TransactionTypeDebit,
	}
}
