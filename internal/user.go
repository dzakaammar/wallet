package internal

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID       string
	Username string
}

func NewUser(username string) (*User, error) {
	if username == "" {
		return nil, WrapErr(ErrInvalidParameter, "invalid username")
	}

	return &User{
		ID:       uuid.NewV4().String(),
		Username: username,
	}, nil
}

type UserRepository interface {
	Store(ctx context.Context, user User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

type UserEvent struct {
	UserID string
}
