package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/domain"
)

func TestNewUser(t *testing.T) {
	t.Run("should create user with email and password", func(t *testing.T) {
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		password, err := domain.NewPassword("ValidPass123©")
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		assert.NotEmpty(t, user.ID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, password, user.Password)
		assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
		assert.Equal(t, user.CreatedAt, user.UpdatedAt)
	})
}

func TestUser_Update(t *testing.T) {
	t.Run("should update user email and password", func(t *testing.T) {
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		password, err := domain.NewPassword("ValidPass123©")
		require.NoError(t, err)

		user := domain.NewUser(email, password)
		originalUser := user

		// Wait a bit to ensure UpdatedAt will be different
		time.Sleep(1 * time.Millisecond)

		newEmail, err := domain.NewEmail("newuser@example.com")
		require.NoError(t, err)

		newPassword, err := domain.NewPassword("NewPass456!%")
		require.NoError(t, err)

		user.Update(newEmail, newPassword)

		assert.Equal(t, originalUser.ID, user.ID)
		assert.Equal(t, originalUser.CreatedAt, user.CreatedAt)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, newPassword, user.Password)
		assert.True(t, user.UpdatedAt.After(originalUser.UpdatedAt))
	})
}
