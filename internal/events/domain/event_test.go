package domain_test

import (
	"testing"
	"time"

	"momento/internal/events/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewEvent(t *testing.T) {
	t.Run("should create event with valid data", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()

		title, err := domain.NewEventTitle("Test event title")
		require.NoError(t, err)

		content, err := domain.NewEventContent("Test event content")
		require.NoError(t, err)

		event := domain.NewEvent(userID, title, content)

		assert.NotEmpty(t, event.ID)
		assert.Equal(t, userID, event.OwnerUserID)
		assert.Equal(t, title, event.Title)
		assert.Equal(t, content, event.Content)
		assert.WithinDuration(t, time.Now().UTC(), event.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), event.UpdatedAt, time.Second)
		assert.Equal(t, event.CreatedAt, event.UpdatedAt)
		assert.Equal(t, time.UTC, event.CreatedAt.Location())
		assert.Equal(t, time.UTC, event.UpdatedAt.Location())
	})
}

func TestEvent_Update(t *testing.T) {
	t.Run("should update event", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		title, _ := domain.NewEventTitle("Test event title")
		content, _ := domain.NewEventContent("Test event content")

		event := domain.NewEvent(userID, title, content)
		newTitle, _ := domain.NewEventTitle("New test event title")
		event.Update(newTitle, content)

		assert.Equal(t, newTitle, event.Title)
		assert.Equal(t, content, event.Content)
		assert.NotEqual(t, event.CreatedAt, event.UpdatedAt)
	})
}
