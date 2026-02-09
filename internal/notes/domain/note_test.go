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
		content, err := domain.NewNoteContent("Test note content")
		require.NoError(t, err)

		note := domain.NewNote(userID, content)

		assert.NotEmpty(t, note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, content, note.Content)
		assert.WithinDuration(t, time.Now().UTC(), note.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), note.UpdatedAt, time.Second)
		assert.Equal(t, note.CreatedAt, note.UpdatedAt)
		assert.Equal(t, time.UTC, note.CreatedAt.Location())
		assert.Equal(t, time.UTC, note.UpdatedAt.Location())
	})
}
