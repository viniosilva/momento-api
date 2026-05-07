package app

import (
	"context"
	"time"

	"momento/internal/auth/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
	FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type JWTService interface {
	Generate(userID string, email domain.Email) (string, error)
}

type SecureTokenService interface {
	Generate(ctx context.Context, userID, email string) (string, error)
	Refresh(ctx context.Context, token string) (userID, email, newToken string, err error)
	Invalidate(ctx context.Context, token string) error
}

type ResetTokenService interface {
	Store(ctx context.Context, token domain.ResetToken, userID string, ttl time.Duration) error
	Validate(ctx context.Context, token domain.ResetToken) (string, error)
	Invalidate(ctx context.Context, token domain.ResetToken) error
}

type EmailSender interface {
	SendResetPasswordEmail(ctx context.Context, to, token string) error
}
