package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrRefreshTokenInvalid  = errors.New("invalid refresh token")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type User struct {
	ID              uuid.UUID
	Email           Email
	Password        Password
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EmailVerifiedAt *time.Time
}

func NewUser(email Email, password Password) User {
	now := time.Now().UTC()

	return User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) UpdatePassword(newPassword Password) {
	u.Password = newPassword
	u.UpdatedAt = time.Now().UTC()
}
