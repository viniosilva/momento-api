package domain

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword       = errors.New("invalid password")
	ErrPasswordTooShort      = errors.New("password must be at least 6 characters")
	ErrPasswordTooLong       = errors.New("password must be less than 64 characters")
	ErrPasswordMissingUpper  = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLower  = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingNumber = errors.New("password must contain at least one number")
	ErrPasswordMissingSymbol = errors.New("password must contain at least one symbol")
	ErrPasswordMismatch      = errors.New("password mismatch")
)

const (
	// minPasswordLength defines minimum password length to ensure basic security
	minPasswordLength = 6
	// maxPasswordLength prevents DoS attacks via bcrypt computation time
	maxPasswordLength = 64
	// bcryptCost balances security and performance (2^12 iterations)
	bcryptCost = 12
)

type Password string

func NewPassword(plainPassword string) (Password, error) {
	if err := ValidatePassword(plainPassword); err != nil {
		return Password(""), err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcryptCost)
	if err != nil {
		return Password(""), fmt.Errorf("failed to hash password: %w", err)
	}

	return Password(string(hashed)), nil
}

func NewPasswordFromHash(hashed string) Password {
	return Password(hashed)
}

func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	var hasUpper, hasLower, hasNumber, hasSymbol bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if !hasUpper {
		return ErrPasswordMissingUpper
	}
	if !hasLower {
		return ErrPasswordMissingLower
	}
	if !hasNumber {
		return ErrPasswordMissingNumber
	}
	if !hasSymbol {
		return ErrPasswordMissingSymbol
	}

	return nil
}

func (p Password) Compare(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(string(p)), []byte(plainPassword))
	if err != nil {
		return ErrPasswordMismatch
	}

	return nil
}
