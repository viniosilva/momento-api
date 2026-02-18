package nethttp_auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyEmail  contextKey = "email"
)

type claims interface {
	GetUserID() string
	GetEmail() string
}

func AuthMiddleware[T claims](jwtService interface {
	Validate(tokenString string) (T, error)
}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				JSON(w, http.StatusUnauthorized, map[string]string{
					"message": "authorization header is required",
				})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
				JSON(w, http.StatusUnauthorized, map[string]string{
					"message": "invalid authorization header format",
				})
				return
			}

			token := parts[1]
			claims, err := jwtService.Validate(token)
			if err != nil {
				JSON(w, http.StatusUnauthorized, map[string]string{
					"message": "invalid or expired token",
				})
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.GetUserID())
			ctx = context.WithValue(ctx, ContextKeyEmail, claims.GetEmail())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
