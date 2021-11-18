package internal

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID       string
	Username string
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.Username, validation.Required),
	)
}

type UserRepository interface {
	Store(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
}
