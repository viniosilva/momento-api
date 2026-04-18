package adapters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"momento/internal/notes/domain"
)

func TestToNoteDocument(t *testing.T) {
	t.Run("should convert domain note to document", func(t *testing.T) {
		id := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)

		note := domain.Note{
			ID:        id.Hex(),
			UserID:    userID.Hex(),
			Title:     domain.NoteTitle("Test Title"),
			Content:   domain.NoteContent("Test Content"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		doc, err := toNoteDocument(note)
		require.NoError(t, err)

		assert.Equal(t, id, doc.ID)
		assert.Equal(t, userID, doc.UserID)
		assert.Equal(t, "Test Title", doc.Title)
		assert.Equal(t, "Test Content", doc.Content)
		assert.Equal(t, now, doc.CreatedAt)
		assert.Equal(t, now, doc.UpdatedAt)
		assert.Nil(t, doc.ArchivedAt)
	})

	t.Run("should return error for invalid note ID", func(t *testing.T) {
		note := domain.Note{
			ID:     "invalid-id",
			UserID: primitive.NewObjectID().Hex(),
		}

		_, err := toNoteDocument(note)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid note ID")
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		note := domain.Note{
			ID:     primitive.NewObjectID().Hex(),
			UserID: "invalid-user-id",
		}

		_, err := toNoteDocument(note)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestToNoteDomain(t *testing.T) {
	t.Run("should convert document to domain note without archived_at", func(t *testing.T) {
		id := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)

		doc := noteDocument{
			ID:        id,
			UserID:    userID,
			Title:     "Test Title",
			Content:   "Test Content",
			CreatedAt: now,
			UpdatedAt: now,
		}

		note := toNoteDomain(doc)

		assert.Equal(t, id.Hex(), note.ID)
		assert.Equal(t, userID.Hex(), note.UserID)
		assert.Equal(t, domain.NoteTitle("Test Title"), note.Title)
		assert.Equal(t, domain.NoteContent("Test Content"), note.Content)
		assert.Equal(t, now, note.CreatedAt)
		assert.Equal(t, now, note.UpdatedAt)
		assert.Nil(t, note.ArchivedAt)
	})

	t.Run("should convert document to domain note with archived_at", func(t *testing.T) {
		id := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now().Truncate(time.Millisecond)
		archivedAt := now.Add(-time.Hour)

		doc := noteDocument{
			ID:         id,
			UserID:     userID,
			Title:      "Archived Note",
			Content:    "Archived Content",
			CreatedAt:  now,
			UpdatedAt:  now,
			ArchivedAt: &archivedAt,
		}

		note := toNoteDomain(doc)

		require.NotNil(t, note.ArchivedAt)
		assert.Equal(t, archivedAt, *note.ArchivedAt)
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
