package rest

import (
	"context"
	"net/http"
	"strings"
	"wallet/internal"
)

type userID string

func authMiddleware(auth *internal.AuthToken) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			header := req.Header.Get("Authorization")
			if header == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			bearer := strings.Split(header, " ")
			if len(bearer) < 2 || bearer[1] == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := auth.ParseToken(bearer[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), new(userID), claims.UserID)

			next.ServeHTTP(w, req.WithContext(ctx))
		}
	}
}
