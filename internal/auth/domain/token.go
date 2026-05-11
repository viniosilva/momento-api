package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var (
	ErrTokenSizeMustBePositive = errors.New("reset token size must be positive")
)

func GenerateSecureToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
