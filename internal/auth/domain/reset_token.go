package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidResetToken = errors.New("invalid reset token")
	ErrExpiredResetToken = errors.New("reset token expired")
)

type ResetToken string

func NewResetToken(size int) (ResetToken, error) {
	if size <= 0 {
		return "", errors.New("token size must be positive")
	}

	token, err := generateSecureToken(size)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return ResetToken(token), nil
}

func (t ResetToken) String() string {
	return string(t)
}

func generateSecureToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
