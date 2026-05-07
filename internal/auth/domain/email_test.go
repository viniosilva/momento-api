package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"momento/internal/auth/domain"
)

func TestNewEmail(t *testing.T) {
	t.Run("should create email with valid format", func(t *testing.T) {
		got, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		assert.Equal(t, "user@example.com", string(got))
	})

	t.Run("should normalize email to lowercase", func(t *testing.T) {
		got, err := domain.NewEmail("USER@EXAMPLE.COM")
		require.NoError(t, err)

		assert.Equal(t, "user@example.com", string(got))
	})

	t.Run("should trim whitespace", func(t *testing.T) {
		got, err := domain.NewEmail("  user@example.com  ")
		require.NoError(t, err)

		assert.Equal(t, "user@example.com", string(got))
	})

	t.Run("should return error for empty email", func(t *testing.T) {
		_, err := domain.NewEmail("")

		assert.ErrorIs(t, err, domain.ErrEmailIsEmpty)
	})

	t.Run("should return error for empty email when has whitespace", func(t *testing.T) {
		_, err := domain.NewEmail(" ")

		assert.ErrorIs(t, err, domain.ErrEmailIsEmpty)
	})

	t.Run("should return error for invalid email format", func(t *testing.T) {
		_, err := domain.NewEmail("invalid-email")

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return error for email without @", func(t *testing.T) {
		_, err := domain.NewEmail("userexample.com")

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return error for email without domain", func(t *testing.T) {
		_, err := domain.NewEmail("user@")

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})
}

func TestEmail_String(t *testing.T) {
	t.Run("should return email as string", func(t *testing.T) {
		emailStr := "user@example.com"
		email, err := domain.NewEmail(emailStr)
		require.NoError(t, err)

		assert.Equal(t, emailStr, email.String())
	})
}
