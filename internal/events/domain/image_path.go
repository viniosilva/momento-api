package domain

import (
	"errors"
)

var (
	ErrInvalidImagePath = errors.New("invalid image path: must match pattern events/{eventID}/{uuid}.{ext}")
)

type ImagePath string

func NewImagePath(value string) (ImagePath, error) {
	if value == "" {
		return ImagePath(""), ErrInvalidImagePath
	}

	return ImagePath(value), nil
}
