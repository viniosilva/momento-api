package domain

import (
	"errors"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	ErrTitleEmpty   = errors.New("title cannot be empty")
	ErrTitleTooLong = errors.New("content exceeds maximum length of 24 characters")
)

type NoteTitle string

func NewNoteTitle(value string) (NoteTitle, error) {
	normalizedValue := strings.TrimSpace(value)

	if normalizedValue == "" {
		return NoteTitle(""), ErrTitleEmpty
	}

	if len(normalizedValue) > 24 {
		return NoteTitle(""), ErrTitleTooLong
	}

	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(normalizedValue)

	return NoteTitle(sanitized), nil
}
