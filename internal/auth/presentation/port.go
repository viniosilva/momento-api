package presentation

import (
	"context"
	"pinnado/internal/auth/application"
)

type AuthService interface {
	Register(ctx context.Context, input application.UserInput) (application.UserOutput, error)
}
