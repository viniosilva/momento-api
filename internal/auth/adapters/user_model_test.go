package adapters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"momento/internal/auth/domain"
)

func TestToUserDocument(t *testing.T) {
	t.Run("should convert domain user to document", func(t *testing.T) {
		id := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)

		user := domain.User{
			ID:        id.Hex(),
			Email:     domain.Email("test@example.com"),
			Password:  domain.Password("hashedpassword"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		doc, err := toUserDocument(user)
		require.NoError(t, err)

		assert.Equal(t, id, doc.ID)
		assert.Equal(t, "test@example.com", doc.Email)
		assert.Equal(t, "hashedpassword", doc.Password)
		assert.Equal(t, now, doc.CreatedAt)
		assert.Equal(t, now, doc.UpdatedAt)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		user := domain.User{
			ID:    "invalid-hex-id",
			Email: domain.Email("test@example.com"),
		}

		_, err := toUserDocument(user)

		assert.Error(t, err)
	})
}

func TestToUserDomain(t *testing.T) {
	t.Run("should convert document to domain user", func(t *testing.T) {
		id := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)

		doc := userDocument{
			ID:        id,
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}

		user := toUserDomain(doc)

		assert.Equal(t, id.Hex(), user.ID)
		assert.Equal(t, domain.Email("test@example.com"), user.Email)
		assert.Equal(t, domain.Password("hashedpassword"), user.Password)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
	})
}
