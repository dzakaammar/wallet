package rest

import (
	"encoding/json"
	"net/http"
	"wallet/internal"
	"wallet/internal/account"
	"wallet/internal/transaction"
	"wallet/internal/user"

	"github.com/gorilla/mux"
)

type Handler struct {
	accountSvc account.Service
	userSvc    user.Service
	trxSvc     transaction.Service
	auth       *internal.AuthToken
}

func NewHandler(accountSvc account.Service, userSvc user.Service, trxSvc transaction.Service, authToken *internal.AuthToken) *Handler {
	return &Handler{
		accountSvc: accountSvc,
		userSvc:    userSvc,
		trxSvc:     trxSvc,
		auth:       authToken,
	}
}

func (h *Handler) Register(r *mux.Router) {
	withToken := authMiddleware(h.auth)
	r.HandleFunc("/balance_read", withToken(h.readBalance)).Methods(http.MethodGet)
	r.HandleFunc("/transfer", withToken(h.transfer)).Methods(http.MethodPost)
	r.HandleFunc("/create_user", h.createUser).Methods(http.MethodPost)
	r.HandleFunc("/top_transactions_per_user", withToken(h.getTopTrxOfUser)).Methods(http.MethodGet)
	r.HandleFunc("/top_users", withToken(h.getTopUserTrxs)).Methods(http.MethodGet)
	r.HandleFunc("/balance_topup", withToken(h.balanceTopUp)).Methods(http.MethodPost)
}

func (h *Handler) readBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDInCtx(ctx)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	balance, err := h.accountSvc.GetUserBalance(ctx, userID)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, struct {
		Balance int `json:"balance"`
	}{
		Balance: balance,
	})
}

func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDInCtx(ctx)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	var reqBody struct {
		ToUsername string `json:"to_username"`
		Amount     int    `json:"amount"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.accountSvc.P2PTransfer(ctx, internal.P2PTransferRequest{
		InitiatorUserID: userID,
		ToUsername:      reqBody.ToUsername,
		Amount:          reqBody.Amount,
	})
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := h.userSvc.Register(r.Context(), reqBody.Username)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	token, err := h.auth.GetTokenFor(userID)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (h *Handler) getTopTrxOfUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDInCtx(ctx)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	res, err := h.trxSvc.GetUserTopTransactions(ctx, userID)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (h *Handler) getTopUserTrxs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := getUserIDInCtx(ctx)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	res, err := h.trxSvc.GetTopUserTransactions(ctx)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, res)
}

func (h *Handler) balanceTopUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDInCtx(ctx)
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	var reqBody struct {
		Amount int `json:"amount"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.accountSvc.TopUp(ctx, internal.TopUpRequest{
		UserID: userID,
		Amount: reqBody.Amount,
	})
	if err != nil {
		mapInternalErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(i)
}
