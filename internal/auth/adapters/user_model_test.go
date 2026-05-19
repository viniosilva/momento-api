package adapters

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"momento/internal/auth/domain"
)

func TestToUserRow(t *testing.T) {
	t.Run("should convert domain user to row", func(t *testing.T) {
		id := uuid.New()
		now := time.Now().Truncate(time.Millisecond)

		user := domain.User{
			ID:        id,
			Email:     domain.Email("test@example.com"),
			Password:  domain.Password("hashedpassword"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		row := toUserRow(user)

		assert.Equal(t, id.String(), row.ID)
		assert.Equal(t, "test@example.com", row.Email)
		assert.Equal(t, "hashedpassword", row.Password)
		assert.Equal(t, now, row.CreatedAt)
		assert.Equal(t, now, row.UpdatedAt)
	})
}

func TestToUserDomain(t *testing.T) {
	t.Run("should convert row to domain user", func(t *testing.T) {
		id := uuid.NewString()
		now := time.Now().Truncate(time.Millisecond)

		row := userRow{
			ID:        id,
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}

		user := toUserDomain(row)

		assert.Equal(t, id, user.ID.String())
		assert.Equal(t, domain.Email("test@example.com"), user.Email)
		assert.Equal(t, domain.Password("hashedpassword"), user.Password)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
	})
}
