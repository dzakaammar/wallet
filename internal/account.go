package internal

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Account struct {
	ID      string
	UserID  string
	Balance int
}

func (a Account) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ID, validation.Required),
		validation.Field(&a.UserID, validation.Required),
		validation.Field(&a.Balance, validation.Min(0)),
	)
}
