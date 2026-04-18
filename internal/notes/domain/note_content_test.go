package domain_test

import (
	"strings"
	"testing"

	"momento/internal/notes/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNoteContent(t *testing.T) {
	t.Run("should create valid note content", func(t *testing.T) {
		content, err := domain.NewNoteContent("Valid note content")
		require.NoError(t, err)

		assert.Equal(t, "Valid note content", string(content))
	})

	t.Run("should trim spaces", func(t *testing.T) {
		content, err := domain.NewNoteContent("  Valid note  ")
		require.NoError(t, err)

		assert.Equal(t, "Valid note", string(content))
	})

	t.Run("should propagate validation errors", func(t *testing.T) {
		_, err := domain.NewNoteContent("")

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should accept exactly 100000 characters", func(t *testing.T) {
		exactContent := strings.Repeat("a", 100000)
		content, err := domain.NewNoteContent(exactContent)
		require.NoError(t, err)

		assert.Equal(t, exactContent, string(content))
	})

	t.Run("should sanitize XSS content", func(t *testing.T) {
		xssContent := "<script>alert('xss')</script>Valid content"
		content, err := domain.NewNoteContent(xssContent)
		require.NoError(t, err)

		assert.Equal(t, "Valid content", string(content))
	})

	t.Run("should sanitize HTML tags", func(t *testing.T) {
		htmlContent := "<p>Valid content</p>"
		content, err := domain.NewNoteContent(htmlContent)
		require.NoError(t, err)

		assert.Equal(t, "Valid content", string(content))
	})

	t.Run("should sanitize iframe tags", func(t *testing.T) {
		iframeContent := "<iframe src='evil.com'></iframe>Valid content"
		content, err := domain.NewNoteContent(iframeContent)
		require.NoError(t, err)

		assert.Equal(t, "Valid content", string(content))
	})

	t.Run("should return error when content is empty", func(t *testing.T) {
		_, err := domain.NewNoteContent("")

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when content is only spaces", func(t *testing.T) {
		_, err := domain.NewNoteContent("   ")

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when content exceeds 100000 characters", func(t *testing.T) {
		longContent := strings.Repeat("a", 100001)
		_, err := domain.NewNoteContent(longContent)

		assert.ErrorIs(t, err, domain.ErrContentTooLong)
	})
}
