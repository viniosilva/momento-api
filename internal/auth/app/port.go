package app

import (
	"context"

	"momento/internal/auth/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
	FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type JWTService interface {
	Generate(userID string, email domain.Email) (string, error)
}

type SecureTokenService interface {
	Generate(ctx context.Context, userID, email string) (string, error)
	Refresh(ctx context.Context, token string) (userID, email, newToken string, err error)
}
