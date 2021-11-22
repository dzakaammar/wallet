package account

import (
	"context"
	"fmt"
	"sync"
	"wallet/internal"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type lock struct {
	lockData map[string]*semaphore.Weighted
	mx       *sync.RWMutex
}

func newLock() *lock {
	return &lock{
		lockData: make(map[string]*semaphore.Weighted),
		mx:       new(sync.RWMutex),
	}
}

func (l *lock) AcquireLock(ctx context.Context, userID string) (func(), error) {
	sem := l.getSemphore(userID)
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return nil, internal.WrapErr(internal.ErrUserTransactionBusy, fmt.Sprintf("user id %s", userID))
	}

	return func() { sem.Release(1) }, nil
}

func (l *lock) AcquireTransferLock(ctx context.Context, debiturUserID, crediturUserID string) (func(), error) {
	wg := new(errgroup.Group)

	var debiturReleaseFn, crediturReleaseFn func()

	wg.Go(func() error {
		var innerErr error
		debiturReleaseFn, innerErr = l.AcquireLock(ctx, debiturUserID)
		return innerErr
	})

	wg.Go(func() error {
		var innerErr error
		crediturReleaseFn, innerErr = l.AcquireLock(ctx, crediturUserID)
		return innerErr
	})

	if err := wg.Wait(); err != nil {
		if debiturReleaseFn != nil {
			debiturReleaseFn()
		}
		if crediturReleaseFn != nil {
			crediturReleaseFn()
		}
		return nil, err
	}

	return func() {
		debiturReleaseFn()
		crediturReleaseFn()
	}, nil
}

func (l *lock) getSemphore(userID string) *semaphore.Weighted {
	l.mx.Lock()
	if _, ok := l.lockData[userID]; !ok {
		l.lockData[userID] = semaphore.NewWeighted(1)
	}
	l.mx.Unlock()

	l.mx.RLock()
	defer l.mx.RUnlock()
	return l.lockData[userID]
}
