package adapters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"momento/internal/events/domain"
)

func TestToEventDocument(t *testing.T) {
	t.Run("should convert domain event to document", func(t *testing.T) {
		id := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)

		event := domain.Event{
			ID:          id.Hex(),
			OwnerUserID: userID.Hex(),
			Title:       domain.EventTitle("Test Title"),
			Content:     domain.EventContent("Test Content"),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		doc, err := toEventDocument(event)
		require.NoError(t, err)

		assert.Equal(t, id, doc.ID)
		assert.Equal(t, userID, doc.OwnerUserID)
		assert.Equal(t, "Test Title", doc.Title)
		assert.Equal(t, "Test Content", doc.Content)
		assert.Equal(t, now, doc.CreatedAt)
		assert.Equal(t, now, doc.UpdatedAt)
		assert.Nil(t, doc.ArchivedAt)
	})

	t.Run("should return error for invalid event ID", func(t *testing.T) {
		event := domain.Event{
			ID:          "invalid-id",
			OwnerUserID: primitive.NewObjectID().Hex(),
		}

		_, err := toEventDocument(event)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		event := domain.Event{
			ID:          primitive.NewObjectID().Hex(),
			OwnerUserID: "invalid-user-id",
		}

		_, err := toEventDocument(event)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestToEventDomain(t *testing.T) {
	t.Run("should convert document to domain event without archived_at", func(t *testing.T) {
		id := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)

		doc := eventDocument{
			ID:          id,
			OwnerUserID: userID,
			Title:       "Test Title",
			Content:     "Test Content",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		event := toEventDomain(doc)

		assert.Equal(t, id.Hex(), event.ID)
		assert.Equal(t, userID.Hex(), event.OwnerUserID)
		assert.Equal(t, domain.EventTitle("Test Title"), event.Title)
		assert.Equal(t, domain.EventContent("Test Content"), event.Content)
		assert.Equal(t, now, event.CreatedAt)
		assert.Equal(t, now, event.UpdatedAt)
		assert.Nil(t, event.ArchivedAt)
	})

	t.Run("should convert document to domain event with archived_at", func(t *testing.T) {
		id := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)
		archivedAt := now.Add(-time.Hour)

		doc := eventDocument{
			ID:          id,
			OwnerUserID: userID,
			Title:       "Archived Event",
			Content:     "Archived Content",
			CreatedAt:   now,
			UpdatedAt:   now,
			ArchivedAt:  &archivedAt,
		}

		event := toEventDomain(doc)

		require.NotNil(t, event.ArchivedAt)
		assert.Equal(t, archivedAt, *event.ArchivedAt)
	})
}

func TestParseObjectID(t *testing.T) {
	t.Run("should parse valid hex string", func(t *testing.T) {
		id := primitive.NewObjectID()

		got, err := parseObjectID(id.Hex())
		require.NoError(t, err)

		assert.Equal(t, id, got)
	})

	t.Run("should return error for invalid hex string", func(t *testing.T) {
		_, err := parseObjectID("invalid-hex")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}
