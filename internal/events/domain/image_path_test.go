package domain_test

import (
	"testing"

	"momento/internal/events/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImagePath(t *testing.T) {
	t.Run("should create valid image path", func(t *testing.T) {
		path, err := domain.NewImagePath("events/abc-123/uuid.jpg")
		require.NoError(t, err)
		assert.Equal(t, domain.ImagePath("events/abc-123/uuid.jpg"), path)
	})

	t.Run("should return error when path is empty", func(t *testing.T) {
		_, err := domain.NewImagePath("")
		assert.ErrorIs(t, err, domain.ErrInvalidImagePath)
	})
}
