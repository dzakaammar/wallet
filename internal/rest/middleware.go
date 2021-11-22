package rest

import (
	"context"
	"net/http"
	"wallet/internal"
)

type ctxKey int

const (
	_ ctxKey = iota
	userIDCtxKey
)

func authMiddleware(auth *internal.AuthToken) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			authHead := req.Header.Get("Authorization")
			if authHead == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// without Bearer as test docs have written
			claims, err := auth.ParseToken(authHead)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := setUserIDInCtx(req.Context(), claims.UserID)
			next.ServeHTTP(w, req.WithContext(ctx))
		}
	}
}

func setUserIDInCtx(parentCtx context.Context, userID string) context.Context {
	return context.WithValue(parentCtx, userIDCtxKey, userID)
}

func getUserIDInCtx(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDCtxKey).(string)
	if !ok || userID == "" {
		return "", internal.ErrDataNotFound
	}
	return userID, nil
}
