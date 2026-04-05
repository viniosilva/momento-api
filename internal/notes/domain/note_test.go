package domain_test

import (
	"testing"
	"time"

	"pinnado/internal/notes/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewNote(t *testing.T) {
	t.Run("should create note with valid data", func(t *testing.T) {
		userID := primitive.NewObjectID()

		title, err := domain.NewNoteTitle("Test note title")
		require.NoError(t, err)

		content, err := domain.NewNoteContent("Test note content")
		require.NoError(t, err)

		note := domain.NewNote(userID, title, content)

		assert.NotEmpty(t, note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, title, note.Title)
		assert.Equal(t, content, note.Content)
		assert.WithinDuration(t, time.Now().UTC(), note.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), note.UpdatedAt, time.Second)
		assert.Equal(t, note.CreatedAt, note.UpdatedAt)
		assert.Equal(t, time.UTC, note.CreatedAt.Location())
		assert.Equal(t, time.UTC, note.UpdatedAt.Location())
	})
}

func TestNote_SetTitle(t *testing.T) {
	t.Run("should set title", func(t *testing.T) {
		userID := primitive.NewObjectID()
		title, _ := domain.NewNoteTitle("Test note title")
		content, _ := domain.NewNoteContent("Test note content")

		note := domain.NewNote(userID, title, content)
		newTitle, _ := domain.NewNoteTitle("New test note title")
		note.SetTitle(newTitle)

		assert.Equal(t, newTitle, note.Title)
	})
}

func TestNote_SetContent(t *testing.T) {
	t.Run("should set content", func(t *testing.T) {
		userID := primitive.NewObjectID()
		title, _ := domain.NewNoteTitle("Test note title")
		content, _ := domain.NewNoteContent("Test note content")

		note := domain.NewNote(userID, title, content)
		newContent, _ := domain.NewNoteContent("New test note content")
		note.SetContent(newContent)

		assert.Equal(t, newContent, note.Content)
	})
}
