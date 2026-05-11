package domain

import (
	"errors"
)

var (
	ErrInvalidResetToken = errors.New("invalid reset token")
	ErrExpiredResetToken = errors.New("reset token expired")
)

type ResetToken string

func NewResetToken(size int) (ResetToken, error) {
	if size <= 0 {
		return "", ErrTokenSizeMustBePositive
	}

	token, err := GenerateSecureToken(size)
	if err != nil {
		return "", err
	}

	return ResetToken(token), nil
}

func (t ResetToken) String() string {
	return string(t)
}
