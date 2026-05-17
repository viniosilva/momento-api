package adapters

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"momento/internal/events/domain"
)

func TestToEventRow(t *testing.T) {
	t.Run("should convert domain event to row", func(t *testing.T) {
		id := uuid.NewString()
		userID := uuid.NewString()
		now := time.Now().Truncate(time.Millisecond)

		event := domain.Event{
			ID:          id,
			OwnerUserID: userID,
			Title:       domain.EventTitle("Test Title"),
			Content:     domain.EventContent("Test Content"),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		row := toEventRow(event)

		assert.Equal(t, id, row.ID)
		assert.Equal(t, userID, row.OwnerUserID)
		assert.Equal(t, "Test Title", row.Title)
		assert.Equal(t, "Test Content", row.Content)
		assert.Equal(t, now, row.CreatedAt)
		assert.Equal(t, now, row.UpdatedAt)
		assert.Nil(t, row.ArchivedAt)
	})
}

func TestToEventDomain(t *testing.T) {
	t.Run("should convert row to domain event without archived_at or images", func(t *testing.T) {
		id := uuid.NewString()
		userID := uuid.NewString()
		now := time.Now().Truncate(time.Millisecond)

		row := eventRow{
			ID:          id,
			OwnerUserID: userID,
			Title:       "Test Title",
			Content:     "Test Content",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		event := toEventDomain(row, nil)

		assert.Equal(t, id, event.ID)
		assert.Equal(t, userID, event.OwnerUserID)
		assert.Equal(t, domain.EventTitle("Test Title"), event.Title)
		assert.Equal(t, domain.EventContent("Test Content"), event.Content)
		assert.Equal(t, now, event.CreatedAt)
		assert.Equal(t, now, event.UpdatedAt)
		assert.Nil(t, event.ArchivedAt)
		assert.Nil(t, event.Metadata)
	})

	t.Run("should convert row to domain event with archived_at", func(t *testing.T) {
		id := uuid.NewString()
		userID := uuid.NewString()
		now := time.Now().Truncate(time.Millisecond)
		archivedAt := now.Add(-time.Hour)

		row := eventRow{
			ID:          id,
			OwnerUserID: userID,
			Title:       "Archived Event",
			Content:     "Archived Content",
			CreatedAt:   now,
			UpdatedAt:   now,
			ArchivedAt:  &archivedAt,
		}

		event := toEventDomain(row, nil)

		require.NotNil(t, event.ArchivedAt)
		assert.Equal(t, archivedAt, *event.ArchivedAt)
	})

	t.Run("should populate metadata from image rows", func(t *testing.T) {
		id := uuid.NewString()
		userID := uuid.NewString()
		now := time.Now().Truncate(time.Millisecond)

		row := eventRow{
			ID:          id,
			OwnerUserID: userID,
			Title:       "Test",
			Content:     "Content",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		imageRows := []eventImageRow{
			{EventID: id, Path: "events/" + id + "/img1.jpg"},
			{EventID: id, Path: "events/" + id + "/img2.jpg"},
		}

		event := toEventDomain(row, imageRows)

		require.NotNil(t, event.Metadata)
		assert.Len(t, event.Metadata.ImagePaths, 2)
	})
}
