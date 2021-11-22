package inmemory

import (
	"context"
	"fmt"
	"sync"
	"wallet/internal"
)

type UserRepository struct {
	indexByUserID   map[string]*internal.User
	indexByUsername map[string]string
	rwLock          *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		indexByUserID:   make(map[string]*internal.User),
		indexByUsername: make(map[string]string),
		rwLock:          new(sync.RWMutex),
	}
}

func (u *UserRepository) Store(ctx context.Context, user internal.User) error {
	u.rwLock.Lock()
	defer u.rwLock.Unlock()

	_, ok := u.indexByUsername[user.Username]
	if ok {
		return internal.WrapErr(internal.ErrDataAlreadyExists, fmt.Sprintf("username : %s", user.Username))
	}

	u.indexByUserID[user.ID] = &user
	u.indexByUsername[user.Username] = user.ID

	return nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (*internal.User, error) {
	u.rwLock.RLock()
	defer u.rwLock.RUnlock()

	userID, ok := u.indexByUsername[username]
	if !ok {
		return nil, internal.WrapErr(internal.ErrDataNotFound, fmt.Sprintf("username : %s", username))
	}

	res := *u.indexByUserID[userID] // return the copy
	return &res, nil
}

func (u *UserRepository) FindByID(ctx context.Context, id string) (*internal.User, error) {
	u.rwLock.RLock()
	defer u.rwLock.RUnlock()

	user, ok := u.indexByUserID[id]
	if !ok {
		return nil, internal.WrapErr(internal.ErrDataNotFound, fmt.Sprintf("id : %s", id))
	}

	res := *user // return the copy
	return &res, nil
}
