package domain

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrEmailIsEmpty = errors.New("email cannot be empty")
var ErrInvalidEmail = errors.New("invalid email format")

type Email string

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	if err := ValidateEmail(value); err != nil {
		return Email(""), err
	}

	return Email(strings.ToLower(value)), nil
}

func ValidateEmail(value string) error {
	if value == "" {
		return ErrEmailIsEmpty
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}

func (e Email) String() string {
	return string(e)
}
