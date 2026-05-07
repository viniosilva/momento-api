package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"momento/internal/auth/domain"
)

func TestNewResetToken(t *testing.T) {
	t.Run("should create a valid reset token", func(t *testing.T) {
		token, err := domain.NewResetToken(32)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, token.String())
	})

	t.Run("should return error when size is zero", func(t *testing.T) {
		_, err := domain.NewResetToken(0)

		assert.Contains(t, err.Error(), "token size must be positive")
	})

	t.Run("should return error when size is negative", func(t *testing.T) {
		_, err := domain.NewResetToken(-1)

		assert.Contains(t, err.Error(), "token size must be positive")
	})

	t.Run("should generate different tokens each time", func(t *testing.T) {
		token1, err := domain.NewResetToken(32)
		require.NoError(t, err)

		token2, err := domain.NewResetToken(32)
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2)
	})

	t.Run("should generate token with correct length for base64 encoding", func(t *testing.T) {
		token, err := domain.NewResetToken(32)
		require.NoError(t, err)

		assert.Len(t, token.String(), 44) // 32 bytes of random data will be 44 characters in base64
	})
}

func TestResetToken_String(t *testing.T) {
	token, err := domain.NewResetToken(16)
	require.NoError(t, err)

	str := token.String()

	assert.Equal(t, string(token), str)
}
