package domain

import (
	"errors"
	"slices"
)

const MaxImages = 100

var (
	ErrMaxImagesReached = errors.New("maximum number of images reached (max 100)")
	ErrImageNotFound    = errors.New("image not found")
)

type EventMetadata struct {
	ImagePaths []ImagePath
}

func NewEventMetadata() EventMetadata {
	return EventMetadata{
		ImagePaths: make([]ImagePath, 0),
	}
}

func (m *EventMetadata) AddImage(path ImagePath) error {
	if m.HasImage(path) {
		return nil
	}

	if len(m.ImagePaths) >= MaxImages {
		return ErrMaxImagesReached
	}

	m.ImagePaths = append(m.ImagePaths, path)

	return nil
}

func (m *EventMetadata) RemoveImage(path ImagePath) error {
	idx := slices.Index(m.ImagePaths, path)
	if idx == -1 {
		return ErrImageNotFound
	}

	m.ImagePaths = slices.Delete(m.ImagePaths, idx, idx+1)

	return nil
}

func (m *EventMetadata) HasImage(path ImagePath) bool {
	return slices.Contains(m.ImagePaths, path)
}
