package domain_test

import (
	"strings"
	"testing"

	"momento/internal/events/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventContent(t *testing.T) {
	t.Run("should create valid event content", func(t *testing.T) {
		content, err := domain.NewEventContent("Valid event content")
		require.NoError(t, err)

		assert.Equal(t, "Valid event content", string(content))
	})

	t.Run("should trim spaces", func(t *testing.T) {
		content, err := domain.NewEventContent("  Valid event  ")
		require.NoError(t, err)

		assert.Equal(t, "Valid event", string(content))
	})

	t.Run("should propagate validation errors", func(t *testing.T) {
		_, err := domain.NewEventContent("")

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should accept exactly 100000 characters", func(t *testing.T) {
		exactContent := strings.Repeat("a", 100000)
		content, err := domain.NewEventContent(exactContent)
		require.NoError(t, err)

		assert.Equal(t, exactContent, string(content))
	})

	t.Run("should sanitize XSS content", func(t *testing.T) {
		xssContent := "<script>alert('xss')</script>Valid content"
		content, err := domain.NewEventContent(xssContent)
		require.NoError(t, err)

		assert.Equal(t, "Valid content", string(content))
	})

	t.Run("should sanitize HTML tags", func(t *testing.T) {
		htmlContent := "<p>Valid content</p>"
		content, err := domain.NewEventContent(htmlContent)
		require.NoError(t, err)

		assert.Equal(t, "Valid content", string(content))
	})

	t.Run("should sanitize iframe tags", func(t *testing.T) {
		iframeContent := "<iframe src='evil.com'></iframe>Valid content"
		content, err := domain.NewEventContent(iframeContent)
		require.NoError(t, err)

		assert.Equal(t, "Valid content", string(content))
	})

	t.Run("should return error when content is empty", func(t *testing.T) {
		_, err := domain.NewEventContent("")

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when content is only spaces", func(t *testing.T) {
		_, err := domain.NewEventContent("   ")

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when content exceeds 100000 characters", func(t *testing.T) {
		longContent := strings.Repeat("a", 100001)
		_, err := domain.NewEventContent(longContent)

		assert.ErrorIs(t, err, domain.ErrContentTooLong)
	})
}