package domain

import (
	"errors"
)

var (
	ErrInvalidVerificationToken = errors.New("invalid verification token")
	ErrExpiredVerificationToken = errors.New("verification token expired")
)

type VerificationToken string

func NewVerificationToken(size int) (VerificationToken, error) {
	if size <= 0 {
		return "", ErrTokenSizeMustBePositive
	}

	token, err := GenerateSecureToken(size)
	if err != nil {
		return "", err
	}

	return VerificationToken(token), nil
}

func (t VerificationToken) String() string {
	return string(t)
}
