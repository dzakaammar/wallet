package transaction

import (
	"context"
	"wallet/internal"
)

type Service interface {
	StoreTransaction(ctx context.Context, data internal.TransferEvent)
	GetUserTopTransactions(ctx context.Context, userID string) ([]*internal.UserTopTransactionsResponse, error)
	GetTopUserTransactions(ctx context.Context) ([]*internal.UserTransactionsResponse, error)
}

type service struct {
	trxRepo  internal.TransactionRepository
	userRepo internal.UserRepository
}

func NewService(trxRepo internal.TransactionRepository, userRepo internal.UserRepository) Service {
	return &service{
		trxRepo:  trxRepo,
		userRepo: userRepo,
	}
}

func (s *service) StoreTransaction(ctx context.Context, data internal.TransferEvent) {
	trxs := data.ToTransactionData()

	for _, trx := range trxs {
		s.trxRepo.Store(ctx, trx)
	}
}

func (s *service) GetUserTopTransactions(ctx context.Context, userID string) ([]*internal.UserTopTransactionsResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data := s.trxRepo.FindTopTransactionsByUserID(ctx, userID, 10)

	res := make([]*internal.UserTopTransactionsResponse, len(data))
	for i, d := range data {
		res[i] = &internal.UserTopTransactionsResponse{
			Username: user.Username,
			Amount:   d.DisplayAmount(),
		}
	}

	return res, nil
}

func (s *service) GetTopUserTransactions(ctx context.Context) ([]*internal.UserTransactionsResponse, error) {
	data := s.trxRepo.FindTopTransactingUser(ctx, 10)

	var res []*internal.UserTransactionsResponse
	for _, d := range data {
		d := d
		user, err := s.userRepo.FindByID(ctx, d.UserID)
		if err != nil {
			continue
		}
		res = append(res, &internal.UserTransactionsResponse{
			Username:        user.Username,
			TransactedValue: d.Amount,
		})
	}

	return res, nil
}
