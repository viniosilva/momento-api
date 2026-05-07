package app

import (
	"time"

	"momento/internal/auth/domain"
)

type UserInput struct {
	Email    string
	Password string
}

type UserOutput struct {
	ID        string
	Email     domain.Email
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token        string
	RefreshToken string
}

type RefreshTokenInput struct {
	RefreshToken string
}

type LogoutInput struct {
	RefreshToken string
}

type ForgotPasswordInput struct {
	Email string
}

type ValidateResetTokenInput struct {
	Token string
}

type ResetPasswordInput struct {
	Token    string
	Password string
}
