package rest

import (
	"context"
	"net/http"
	"wallet/internal"

	"github.com/gorilla/mux"
)

type AccountService interface {
	GetBalance(ctx context.Context, userID string) (balance int64, err error)
	TopUp(ctx context.Context, userID string, amount string) error
}

type UserService interface {
	Register(ctx context.Context, username string) (id string, err error)
}

type TransactionService interface {
	Transfer(ctx context.Context, param internal.TransferRequest) error
	GetTopTransactionsPerUser(ctx context.Context, userID string) ([]*internal.TransactionPerUserResponse, error)
	GetTopUserTransactions(ctx context.Context) ([]*internal.UserTransactionsResponse, error)
}

type Handler struct {
	accountSvc AccountService
	userSvc    UserService
	trxSvc     TransactionService
}

func NewHandler(accountSvc AccountService, userSvc UserService, trxSvc TransactionService) *Handler {
	return &Handler{
		accountSvc: accountSvc,
		userSvc:    userSvc,
		trxSvc:     trxSvc,
	}
}

func (h *Handler) Register(auth *internal.AuthToken, r *mux.Router) {
	protect := authMiddleware(auth)
	r.HandleFunc("/balance_read", protect(h.readBalance)).Methods(http.MethodGet)
	r.HandleFunc("/transfer", protect(h.transfer)).Methods(http.MethodPost)
	r.HandleFunc("/create_user", h.createUser).Methods(http.MethodPost)
	r.HandleFunc("/top_transactions_per_user", protect(h.topTrxOfUser)).Methods(http.MethodGet)
	r.HandleFunc("/top_users", protect(h.topUserTrxs)).Methods(http.MethodGet)
	r.HandleFunc("/balance_topup", protect(h.balanceTopUp)).Methods(http.MethodPost)
}

func (h *Handler) readBalance(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) topTrxOfUser(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) topUserTrxs(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) balanceTopUp(w http.ResponseWriter, r *http.Request) {}
