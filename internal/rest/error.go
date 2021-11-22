package rest

import (
	"errors"
	"net/http"
	"wallet/internal"
)

func mapInternalErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, internal.ErrDataNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, internal.ErrDataAlreadyExists):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, internal.ErrDataNotFound):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, internal.ErrInvalidParameter):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, internal.ErrUserTransactionBusy):
		w.WriteHeader(http.StatusLocked)
	case errors.Is(err, internal.ErrInsufficientBalance):
		w.WriteHeader(http.StatusExpectationFailed)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
