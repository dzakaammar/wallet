package user

import (
	"context"
	"wallet/internal"
	"wallet/internal/account"
)

type EventHandler interface {
	UserWasCreated(internal.UserEvent)
}

type Service interface {
	Register(ctx context.Context, username string) (id string, err error)
}

func NewService(userRepo internal.UserRepository, userEventHandler EventHandler) Service {
	return &service{
		userRepo:        userRepo,
		userEvenHandler: userEventHandler,
	}
}

type service struct {
	userRepo        internal.UserRepository
	userEvenHandler EventHandler
}

func (s *service) Register(ctx context.Context, username string) (string, error) {
	user, err := internal.NewUser(username)
	if err != nil {
		return "", err
	}

	err = s.userRepo.Store(ctx, *user)
	if err != nil {
		return "", err
	}

	s.userEvenHandler.UserWasCreated(internal.UserEvent{
		UserID: user.ID,
	})

	return user.ID, nil
}

type userEventHandler struct {
	accountSvc account.Service
}

func NewEventHandler(accountSvc account.Service) EventHandler {
	return &userEventHandler{
		accountSvc: accountSvc,
	}
}

func (u *userEventHandler) UserWasCreated(event internal.UserEvent) {
	_ = u.accountSvc.CreateAccount(context.Background(), event.UserID)
}
