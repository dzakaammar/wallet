package internal

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type AccountRepository interface {
	Store(ctx context.Context, account Account) error
	FindByUserID(ctx context.Context, userID string) (*Account, error)
	Updates(cxt context.Context, account ...Account) error
}

type Account struct {
	ID      string
	UserID  string
	Balance int
}

func NewAccount(userID string) (*Account, error) {
	if userID == "" {
		return nil, WrapErr(ErrInvalidParameter, "invalid userID")
	}

	return &Account{
		ID:      uuid.NewV4().String(),
		UserID:  userID,
		Balance: 0,
	}, nil
}

func (a *Account) TransferTo(creditur *Account, amount int) error {
	if a.ID == creditur.ID {
		return WrapErr(ErrInvalidParameter, "cannot transfer to the same account")
	}

	if eligible := a.isEligibleToWithdraw(amount); !eligible {
		return ErrInsufficientBalance
	}

	a.Balance -= amount
	creditur.Balance += amount

	return nil
}

func (a *Account) Deposit(amount int) {
	a.Balance += amount
}

func (a *Account) isEligibleToWithdraw(amount int) bool {
	return a.Balance-amount > 0
}
