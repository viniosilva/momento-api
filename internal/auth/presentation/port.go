package presentation

import (
	"context"
	"momento/internal/auth/application"
)

type AuthService interface {
	Register(ctx context.Context, input application.UserInput) (application.UserOutput, error)
	Login(ctx context.Context, input application.LoginInput) (application.LoginOutput, error)
}
