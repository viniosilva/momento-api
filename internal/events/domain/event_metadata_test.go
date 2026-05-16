package domain_test

import (
	"testing"

	"momento/internal/events/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventMetadata(t *testing.T) {
	t.Run("should create empty metadata", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		assert.NotNil(t, metadata.ImagePaths)
		assert.Empty(t, metadata.ImagePaths)
	})
}

func TestEventMetadata_AddImage(t *testing.T) {
	t.Run("should add image", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		path, err := domain.NewImagePath("events/event-id/img.jpg")
		require.NoError(t, err)

		err = metadata.AddImage(path)
		assert.NoError(t, err)
		assert.Len(t, metadata.ImagePaths, 1)
		assert.True(t, metadata.HasImage(path))
	})

	t.Run("should not duplicate image", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		path, _ := domain.NewImagePath("events/event-id/img.jpg")

		err := metadata.AddImage(path)
		require.NoError(t, err)

		err = metadata.AddImage(path)
		assert.NoError(t, err)
		assert.Len(t, metadata.ImagePaths, 1)
	})

	t.Run("should return error when max images reached", func(t *testing.T) {
		metadata := domain.NewEventMetadata()

		for i := 0; i < domain.MaxImages; i++ {
			p, _ := domain.NewImagePath("events/event-id/img.jpg")
			p = domain.ImagePath(string(p) + string(rune('0'+i)))
			err := metadata.AddImage(p)
			require.NoError(t, err)
		}

		extra, _ := domain.NewImagePath("events/event-id/extra.jpg")
		err := metadata.AddImage(extra)
		assert.ErrorIs(t, err, domain.ErrMaxImagesReached)
		assert.Len(t, metadata.ImagePaths, domain.MaxImages)
	})
}

func TestEventMetadata_RemoveImage(t *testing.T) {
	t.Run("should remove existing image", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		path, _ := domain.NewImagePath("events/event-id/img.jpg")

		err := metadata.AddImage(path)
		require.NoError(t, err)

		err = metadata.RemoveImage(path)
		assert.NoError(t, err)
		assert.Empty(t, metadata.ImagePaths)
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		path, _ := domain.NewImagePath("events/event-id/img.jpg")

		err := metadata.RemoveImage(path)
		assert.ErrorIs(t, err, domain.ErrImageNotFound)
	})
}

func TestEventMetadata_HasImage(t *testing.T) {
	t.Run("should return true when image exists", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		path, _ := domain.NewImagePath("events/event-id/img.jpg")
		metadata.AddImage(path)

		assert.True(t, metadata.HasImage(path))
	})

	t.Run("should return false when image does not exist", func(t *testing.T) {
		metadata := domain.NewEventMetadata()
		path, _ := domain.NewImagePath("events/event-id/img.jpg")

		assert.False(t, metadata.HasImage(path))
	})
}
