package internal

import (
	"context"
)

type UserTopTransactionsResponse struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type UserTransactionsResponse struct {
	Username        string `json:"username"`
	TransactedValue int    `json:"transacted_value"`
}

type TransactionType int

const (
	_ TransactionType = iota
	DebitTransaction
	CreditTransaction
)

type TransactionData struct {
	UserID string
	Amount int
	Type   TransactionType
}

func (t TransactionData) DisplayAmount() int {
	if t.Type == DebitTransaction {
		return t.Amount * -1
	}

	return t.Amount
}

type TransactionsData []TransactionData

// to satisfy sort.Interface
func (t TransactionsData) Len() int           { return len(t) }
func (t TransactionsData) Less(i, j int) bool { return t[i].Amount < t[j].Amount }
func (t TransactionsData) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type TransactionRepository interface {
	Store(context.Context, TransactionData) error
	FindTopTransactionsByUserID(ctx context.Context, userID string, count int) TransactionsData
	FindTopTransactingUser(ctx context.Context, count int) TransactionsData
}
