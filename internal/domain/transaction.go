package domain

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          string
	Amount      float64
	Type        TransactionType
	Description string
	Date        time.Time
}
