package ports

import (
	"context"

	"momento/internal/auth/app"
)

type AuthService interface {
	Register(ctx context.Context, input app.UserInput) (app.UserOutput, error)
	Login(ctx context.Context, input app.LoginInput) (app.LoginOutput, error)
	RefreshToken(ctx context.Context, input app.RefreshTokenInput) (app.LoginOutput, error)
	Logout(ctx context.Context, input app.LogoutInput) error
	ForgotPassword(ctx context.Context, input app.ForgotPasswordInput) error
	ResetPassword(ctx context.Context, input app.ResetPasswordInput) error
	ValidateResetToken(ctx context.Context, input app.ValidateResetTokenInput) (string, error)
}
