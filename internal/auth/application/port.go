package application

import (
	"context"

	"pinnado/internal/auth/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}
