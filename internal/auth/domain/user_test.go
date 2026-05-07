package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"momento/internal/auth/domain"
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
		assert.Nil(t, user.EmailVerifiedAt)
		assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
		assert.Equal(t, user.CreatedAt, user.UpdatedAt)
	})
}

func TestUser_UpdatePassword(t *testing.T) {
	t.Run("should update user password and updatedAt", func(t *testing.T) {
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		oldPassword, err := domain.NewPassword("oldPass123*")
		require.NoError(t, err)

		user := domain.NewUser(email, oldPassword)

		password, err := domain.NewPassword("ValidPass123%")
		require.NoError(t, err)

		oldUpdatedAt := user.UpdatedAt
		user.UpdatePassword(password)

		assert.Equal(t, password, user.Password)
		assert.NotEqual(t, oldUpdatedAt, user.UpdatedAt)
	})
}
