package nethttp

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyEmail  contextKey = "email"
)

func AuthMiddleware(jwtService JWTService) func(http.Handler) http.Handler {
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

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
