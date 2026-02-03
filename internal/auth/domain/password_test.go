package domain_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/domain"
)

func TestNewPassword(t *testing.T) {
	t.Run("should create password with valid complexity", func(t *testing.T) {
		password, err := domain.NewPassword("ValidPass123©®€!")
		require.NoError(t, err)

		assert.NotEmpty(t, string(password))
		assert.NotEqual(t, "ValidPass123©®€!", string(password))
	})

	t.Run("should create password with simple symbols", func(t *testing.T) {
		password, err := domain.NewPassword("ValidPass123!@#")
		require.NoError(t, err)

		assert.NotEmpty(t, string(password))
		assert.NotEqual(t, "ValidPass123!@#", string(password))
	})

	t.Run("should accept password at maximum length", func(t *testing.T) {
		longPassword := strings.Repeat("A", 58) + "a1©"
		password, err := domain.NewPassword(longPassword)

		require.NoError(t, err)
		assert.NotEmpty(t, string(password))
	})

	t.Run("should return error for password too short", func(t *testing.T) {
		_, err := domain.NewPassword("Sh1©")

		assert.ErrorIs(t, err, domain.ErrPasswordTooShort)

	})

	t.Run("should return error for password too long", func(t *testing.T) {
		longPassword := strings.Repeat("A", 61) + "a1©"
		_, err := domain.NewPassword(longPassword)

		assert.ErrorIs(t, err, domain.ErrPasswordTooLong)
	})

	t.Run("should return error for password without uppercase", func(t *testing.T) {
		_, err := domain.NewPassword("lowercase123!")

		assert.ErrorIs(t, err, domain.ErrPasswordMissingUpper)
	})

	t.Run("should return error for password without lowercase", func(t *testing.T) {
		_, err := domain.NewPassword("UPPERCASE123!")

		assert.ErrorIs(t, err, domain.ErrPasswordMissingLower)
	})

	t.Run("should return error for password without number", func(t *testing.T) {
		_, err := domain.NewPassword("NoNumberPass!")

		assert.ErrorIs(t, err, domain.ErrPasswordMissingNumber)
	})

	t.Run("should return error for password without symbol", func(t *testing.T) {
		_, err := domain.NewPassword("NoSymbolPass123")

		assert.ErrorIs(t, err, domain.ErrPasswordMissingSymbol)
	})
}

func TestPassword_Compare(t *testing.T) {
	t.Run("should return nil for matching password", func(t *testing.T) {
		plainPassword := "ValidPass123©"
		password, err := domain.NewPassword(plainPassword)
		require.NoError(t, err)

		err = password.Compare(plainPassword)
		assert.NoError(t, err)
	})

	t.Run("should return error for non-matching password", func(t *testing.T) {
		password, err := domain.NewPassword("ValidPass123©")
		require.NoError(t, err)

		err = password.Compare("WrongPassword123©")

		assert.ErrorIs(t, err, domain.ErrPasswordMismatch)
	})

	t.Run("should return error for empty password comparison", func(t *testing.T) {
		password, err := domain.NewPassword("ValidPass123©")
		require.NoError(t, err)

		err = password.Compare("")
		assert.ErrorIs(t, err, domain.ErrPasswordMismatch)
	})
}

func TestNewPasswordFromHash(t *testing.T) {
	t.Run("should create password from hash", func(t *testing.T) {
		originalPassword, err := domain.NewPassword("ValidPass123©")
		require.NoError(t, err)

		hash := string(originalPassword)
		passwordFromHash := domain.NewPasswordFromHash(hash)

		assert.Equal(t, hash, string(passwordFromHash))
		assert.NoError(t, passwordFromHash.Compare("ValidPass123©"))
	})

	t.Run("should compare correctly after creating from hash", func(t *testing.T) {
		originalPassword, err := domain.NewPassword("TestPass456©")
		require.NoError(t, err)

		hash := string(originalPassword)
		passwordFromHash := domain.NewPasswordFromHash(hash)

		err = passwordFromHash.Compare("TestPass456©")
		assert.NoError(t, err)

		err = passwordFromHash.Compare("WrongPass456©")
		assert.ErrorIs(t, err, domain.ErrPasswordMismatch)
	})
}
