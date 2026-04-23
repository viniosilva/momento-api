package domain_test

import (
	"testing"

	"momento/internal/events/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventTitle(t *testing.T) {
	t.Run("should create a valid EventTitle", func(t *testing.T) {
		title, err := domain.NewEventTitle("Title")
		require.NoError(t, err)

		assert.Equal(t, "Title", string(title))
	})

	t.Run("should return error when title is empty", func(t *testing.T) {
		_, err := domain.NewEventTitle("")
		assert.ErrorIs(t, err, domain.ErrTitleEmpty)
	})

	t.Run("should return error when title exceeds maximum length", func(t *testing.T) {
		longTitle := "This title is way too long to be valid"
		_, err := domain.NewEventTitle(longTitle)
		assert.ErrorIs(t, err, domain.ErrTitleTooLong)
	})

	t.Run("should trim whitespace from title", func(t *testing.T) {
		title, err := domain.NewEventTitle("  Title with spaces  ")
		require.NoError(t, err)

		assert.Equal(t, "Title with spaces", string(title))
	})
}