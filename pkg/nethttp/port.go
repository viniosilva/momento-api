package nethttp

import "pinnado/internal/auth/infrastructure"

type JWTService interface {
	Validate(tokenString string) (*infrastructure.Claims, error)
}
