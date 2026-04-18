package domain_test

import (
	"testing"

	"momento/internal/notes/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNoteTitle(t *testing.T) {
	t.Run("should create a valid NoteTitle", func(t *testing.T) {
		title, err := domain.NewNoteTitle("Title")
		require.NoError(t, err)

		assert.Equal(t, "Title", string(title))
	})

	t.Run("should return error when title is empty", func(t *testing.T) {
		_, err := domain.NewNoteTitle("")
		assert.ErrorIs(t, err, domain.ErrTitleEmpty)
	})

	t.Run("should return error when title exceeds maximum length", func(t *testing.T) {
		longTitle := "This title is way too long to be valid"
		_, err := domain.NewNoteTitle(longTitle)
		assert.ErrorIs(t, err, domain.ErrTitleTooLong)
	})

	t.Run("should trim whitespace from title", func(t *testing.T) {
		title, err := domain.NewNoteTitle("  Title with spaces  ")
		require.NoError(t, err)

		assert.Equal(t, "Title with spaces", string(title))
	})
}
