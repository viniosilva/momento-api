package infrastructure_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/domain"
	"pinnado/internal/auth/infrastructure"
)

func TestNewJWTService(t *testing.T) {
	t.Run("should create jwt service", func(t *testing.T) {
		service := infrastructure.NewJWTService("test-secret", 24*time.Hour)

		assert.NotNil(t, service)
	})
}

func TestJWTService_Generate(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24 * time.Hour
	service := infrastructure.NewJWTService(secret, expiration)

	t.Run("should generate valid token", func(t *testing.T) {
		userID := "507f1f77bcf86cd799439011"
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		token, err := service.Generate(userID, email)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		email1, err := domain.NewEmail("user1@example.com")
		require.NoError(t, err)

		email2, err := domain.NewEmail("user2@example.com")
		require.NoError(t, err)

		token1, err := service.Generate("user1", email1)
		require.NoError(t, err)

		token2, err := service.Generate("user2", email2)
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2)
	})
}

func TestJWTService_Validate(t *testing.T) {
	secret := "test-secret-key"
	expiration := 5 * time.Second
	service := infrastructure.NewJWTService(secret, expiration)

	t.Run("should validate valid token", func(t *testing.T) {
		userID := "507f1f77bcf86cd799439011"
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		token, err := service.Generate(userID, email)
		require.NoError(t, err)

		claims, err := service.Validate(token)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.GetUserID())
		assert.Equal(t, string(email), claims.GetEmail())
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.IssuedAt)
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		_, err := service.Validate(invalidToken)

		assert.ErrorIs(t, err, infrastructure.ErrInvalidToken)
	})

	t.Run("should return error for token with wrong secret", func(t *testing.T) {
		otherService := infrastructure.NewJWTService("different-secret", expiration)
		email, err := domain.NewEmail("user1@example.com")
		require.NoError(t, err)

		token, err := otherService.Generate("user1", email)
		require.NoError(t, err)

		_, err = service.Validate(token)

		assert.ErrorIs(t, err, infrastructure.ErrInvalidToken)
	})

	t.Run("should return error for expired token", func(t *testing.T) {
		expiredService := infrastructure.NewJWTService(secret, -5*time.Second)
		email, err := domain.NewEmail("user1@example.com")
		require.NoError(t, err)

		token, err := expiredService.Generate("user1", email)
		require.NoError(t, err)

		_, err = service.Validate(token)

		assert.ErrorIs(t, err, infrastructure.ErrExpiredToken)
	})
}
