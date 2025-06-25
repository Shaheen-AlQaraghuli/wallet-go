package types

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"
)

func GetTransactionTypes() []TransactionType {
	return []TransactionType{
		TransactionTypeCredit,
		TransactionTypeDebit,
	}
}
