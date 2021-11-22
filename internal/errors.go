package internal

import (
	"errors"
	"fmt"
)

var (
	ErrDataAlreadyExists   = errors.New("data already exists")
	ErrDataNotFound        = errors.New("data not found")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrUserTransactionBusy = errors.New("user transaction is busy")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Error struct {
	err error
	msg string
}

func WrapErr(err error, msg string) *Error {
	return &Error{
		err: err,
		msg: msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s. reason: %s", e.err.Error(), e.msg)
}

func (e *Error) Unwrap() error {
	return e.err
}
