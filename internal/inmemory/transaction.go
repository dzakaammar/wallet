package inmemory

import (
	"context"
	"sort"
	"sync"
	"wallet/internal"
)

type sumBucket struct {
	data map[string]int
	lock *sync.RWMutex
}

func (s *sumBucket) Add(userID string, amount int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if sum, ok := s.data[userID]; ok {
		s.data[userID] = sum + amount
		return
	}
	s.data[userID] = amount
}

func (r *sumBucket) Get() internal.TransactionsData {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var res internal.TransactionsData
	for userID, amount := range r.data {
		res = append(res, internal.TransactionData{
			UserID: userID,
			Amount: amount,
			Type:   internal.DebitTransaction,
		})
	}

	return res
}

type recordBucket struct {
	data map[string]internal.TransactionsData
	lock *sync.RWMutex
}

func (r *recordBucket) Add(data internal.TransactionData) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.data[data.UserID]; ok {
		r.data[data.UserID] = append(r.data[data.UserID], data)
		return
	}
	r.data[data.UserID] = []internal.TransactionData{data}
}

func (r *recordBucket) Get(userID string) internal.TransactionsData {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if _, ok := r.data[userID]; !ok {
		return nil
	}
	return r.data[userID]
}

type TransactionRepository struct {
	debitedAmountOfUser *sumBucket
	recordBucket        *recordBucket
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		debitedAmountOfUser: &sumBucket{
			data: make(map[string]int),
			lock: new(sync.RWMutex),
		},
		recordBucket: &recordBucket{
			data: make(map[string]internal.TransactionsData),
			lock: new(sync.RWMutex),
		},
	}
}

func (t *TransactionRepository) Store(ctx context.Context, data internal.TransactionData) error {
	if data.Type == internal.DebitTransaction {
		t.debitedAmountOfUser.Add(data.UserID, data.Amount)
	}

	t.recordBucket.Add(data)
	return nil
}

func (t *TransactionRepository) FindTopTransactionsByUserID(ctx context.Context, userID string, count int) internal.TransactionsData {
	list := t.recordBucket.Get(userID)

	sort.Sort(sort.Reverse(list))

	if len(list) <= count {
		return list
	}

	return list[:count-1]
}

func (t *TransactionRepository) FindTopTransactingUser(ctx context.Context, count int) internal.TransactionsData {
	list := t.debitedAmountOfUser.Get()
	sort.Sort(sort.Reverse(list))

	if len(list) <= count {
		return list
	}

	return list[:count-1]
}
