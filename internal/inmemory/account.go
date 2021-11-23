package inmemory

import (
	"context"
	"fmt"
	"sync"
	"wallet/internal"
)

var _ internal.AccountRepository = &AccountRepository{}

type AccountRepository struct {
	// map id to account
	data map[string]*internal.Account
	// map userid to account id
	indexByUserID map[string]string

	lock *sync.RWMutex
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		data:          make(map[string]*internal.Account),
		indexByUserID: make(map[string]string),
		lock:          new(sync.RWMutex),
	}
}

func (a *AccountRepository) Store(ctx context.Context, account internal.Account) error {
	a.lock.RLock()
	if _, ok := a.indexByUserID[account.UserID]; ok {
		a.lock.RUnlock()
		return internal.WrapErr(internal.ErrDataAlreadyExists, "user account already exists")
	}

	if _, ok := a.data[account.ID]; ok {
		a.lock.RUnlock()
		return internal.WrapErr(internal.ErrDataAlreadyExists, "duplicate account id")
	}
	a.lock.RUnlock()

	a.lock.Lock()
	a.data[account.ID] = &account
	a.indexByUserID[account.UserID] = account.ID
	a.lock.Unlock()

	return nil
}

func (a *AccountRepository) FindByUserID(ctx context.Context, userID string) (*internal.Account, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if accountID, ok := a.indexByUserID[userID]; !ok {
		return nil, internal.WrapErr(internal.ErrDataNotFound, fmt.Sprintf("account with userID %s is not found", userID))
	} else {
		acc := *a.data[accountID] // return the copy
		return &acc, nil
	}
}

func (a *AccountRepository) Updates(ctx context.Context, accounts ...internal.Account) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, account := range accounts {
		account := account
		if _, ok := a.data[account.ID]; !ok {
			return internal.WrapErr(internal.ErrDataNotFound, fmt.Sprintf("account id %s", account.ID))
		}
	}

	for _, account := range accounts {
		account := account
		a.data[account.ID] = &account
	}

	return nil
}
