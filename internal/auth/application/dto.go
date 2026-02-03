package application

import (
	"pinnado/internal/auth/domain"
	"time"
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
