package account

import (
	"context"
	"wallet/internal"
	"wallet/internal/transaction"
)

type EventHandler interface {
	TransferWasSucceed(internal.TransferEvent)
}

type Locker interface {
	AcquireLock(ctx context.Context, userID string) (func(), error)
	AcquireTransferLock(ctx context.Context, debiturUserID, crediturUserID string) (func(), error)
}

type Service interface {
	P2PTransfer(ctx context.Context, param internal.P2PTransferRequest) error
	TopUp(ctx context.Context, param internal.TopUpRequest) error
	GetUserBalance(ctx context.Context, userID string) (balance int, err error)
	CreateAccount(ctx context.Context, userID string) error
}

type service struct {
	accountRepo  internal.AccountRepository
	userRepo     internal.UserRepository
	eventHandler EventHandler
	locker       Locker
}

func NewService(
	accountRepo internal.AccountRepository,
	userRepo internal.UserRepository,
	eventHandler EventHandler,
	locker Locker,
) Service {
	return &service{
		accountRepo:  accountRepo,
		userRepo:     userRepo,
		eventHandler: eventHandler,
		locker:       locker,
	}
}

func (s *service) GetUserBalance(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, internal.WrapErr(internal.ErrInvalidParameter, "invalid userID")
	}

	account, err := s.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

func (s *service) CreateAccount(ctx context.Context, userID string) error {
	account, err := internal.NewAccount(userID)
	if err != nil {
		return err
	}

	err = s.accountRepo.Store(ctx, *account)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) TopUp(ctx context.Context, req internal.TopUpRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	releaseLockFn, err := s.locker.AcquireLock(ctx, req.UserID)
	if err != nil {
		return nil
	}
	defer releaseLockFn()

	account, err := s.accountRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	account.Deposit(req.Amount)
	return s.accountRepo.Updates(ctx, *account)
}

func (s *service) P2PTransfer(ctx context.Context, req internal.P2PTransferRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	creditur, err := s.userRepo.FindByUsername(ctx, req.ToUsername)
	if err != nil {
		return err
	}

	releaseLockFn, err := s.locker.AcquireTransferLock(ctx, req.InitiatorUserID, creditur.ID)
	if err != nil {
		return nil
	}
	defer releaseLockFn()

	debiturAccount, err := s.accountRepo.FindByUserID(ctx, req.InitiatorUserID)
	if err != nil {
		return err
	}

	crediturAccount, err := s.accountRepo.FindByUserID(ctx, creditur.ID)
	if err != nil {
		return err
	}

	err = debiturAccount.TransferTo(crediturAccount, req.Amount)
	if err != nil {
		return err
	}

	err = s.accountRepo.Updates(ctx, *debiturAccount, *crediturAccount)
	if err != nil {
		return err
	}

	s.eventHandler.TransferWasSucceed(internal.TransferEvent{
		DebitUserID:  debiturAccount.UserID,
		CreditUserID: crediturAccount.UserID,
		Amount:       req.Amount,
	})

	return nil
}

type eventHandler struct {
	trxSvc transaction.Service
}

func NewEventHandler(trxSvc transaction.Service) EventHandler {
	return &eventHandler{
		trxSvc: trxSvc,
	}
}

func (e *eventHandler) TransferWasSucceed(event internal.TransferEvent) {
	e.trxSvc.StoreTransaction(context.Background(), event)
}
