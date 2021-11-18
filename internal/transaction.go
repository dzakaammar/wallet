package internal

import "time"

type Transaction struct {
	ID           string
	Amount       int
	DebitUserID  string
	CreditUserID string
	CreatedAt    time.Time
}
