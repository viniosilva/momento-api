package domain_test

import (
	"testing"
	"time"

	"momento/internal/events/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/google/uuid"
)

func TestNewEvent(t *testing.T) {
	t.Run("should create event with valid data", func(t *testing.T) {
		userID := uuid.NewString()

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

func TestEvent_AddImage(t *testing.T) {
	t.Run("should add image to event", func(t *testing.T) {
		userID := uuid.NewString()
		title, _ := domain.NewEventTitle("Title")
		content, _ := domain.NewEventContent("Content")
		event := domain.NewEvent(userID, title, content)

		path, _ := domain.NewImagePath("events/event-id/img.jpg")
		err := event.AddImage(path)

		assert.NoError(t, err)
		require.NotNil(t, event.Metadata)
		assert.True(t, event.Metadata.HasImage(path))
		assert.WithinDuration(t, time.Now().UTC(), event.UpdatedAt, time.Second)
	})

	t.Run("should initialize metadata when nil", func(t *testing.T) {
		event := domain.Event{
			ID:    uuid.NewString(),
			Title: "Title",
		}
		assert.Nil(t, event.Metadata)

		path, _ := domain.NewImagePath("events/event-id/img.jpg")
		err := event.AddImage(path)

		assert.NoError(t, err)
		require.NotNil(t, event.Metadata)
		assert.True(t, event.Metadata.HasImage(path))
	})

	t.Run("should return error when max images reached", func(t *testing.T) {
		userID := uuid.NewString()
		title, _ := domain.NewEventTitle("Title")
		content, _ := domain.NewEventContent("Content")
		event := domain.NewEvent(userID, title, content)

		for i := 0; i < domain.MaxImages; i++ {
			p := domain.ImagePath(string(rune('0'+i)))
			event.AddImage(p)
		}

		extra, _ := domain.NewImagePath("events/event-id/extra.jpg")
		err := event.AddImage(extra)
		assert.ErrorIs(t, err, domain.ErrMaxImagesReached)
	})
}

func TestEvent_RemoveImage(t *testing.T) {
	t.Run("should remove existing image", func(t *testing.T) {
		userID := uuid.NewString()
		title, _ := domain.NewEventTitle("Title")
		content, _ := domain.NewEventContent("Content")
		event := domain.NewEvent(userID, title, content)

		path, _ := domain.NewImagePath("events/event-id/img.jpg")
		event.AddImage(path)

		err := event.RemoveImage(path)
		assert.NoError(t, err)
		require.NotNil(t, event.Metadata)
		assert.False(t, event.Metadata.HasImage(path))
		assert.WithinDuration(t, time.Now().UTC(), event.UpdatedAt, time.Second)
	})

	t.Run("should return error when metadata is nil", func(t *testing.T) {
		event := domain.Event{
			ID:    uuid.NewString(),
			Title: "Title",
		}
		assert.Nil(t, event.Metadata)

		path, _ := domain.NewImagePath("events/event-id/img.jpg")
		err := event.RemoveImage(path)
		assert.ErrorIs(t, err, domain.ErrImageNotFound)
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		userID := uuid.NewString()
		title, _ := domain.NewEventTitle("Title")
		content, _ := domain.NewEventContent("Content")
		event := domain.NewEvent(userID, title, content)

		path, _ := domain.NewImagePath("events/event-id/img.jpg")
		err := event.RemoveImage(path)
		assert.ErrorIs(t, err, domain.ErrImageNotFound)
	})
}

func TestEvent_Update(t *testing.T) {
	t.Run("should update event", func(t *testing.T) {
		userID := uuid.NewString()
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
