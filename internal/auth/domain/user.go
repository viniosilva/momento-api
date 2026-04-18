package domain

import (
	"errors"
	"time"

	"momento/pkg/uid"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID        string
	Email     Email
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email Email, password Password) User {
	now := time.Now().UTC()

	return User{
		ID:        uid.New(),
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
