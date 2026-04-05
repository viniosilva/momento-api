package domain

import (
	"errors"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	ErrContentEmpty   = errors.New("content cannot be empty")
	ErrContentTooLong = errors.New("content exceeds maximum length of 100000 characters")
)

type NoteContent string

func NewNoteContent(value string) (NoteContent, error) {
	normalizedValue := strings.TrimSpace(value)

	if normalizedValue == "" {
		return NoteContent(""), ErrContentEmpty
	}

	if len(normalizedValue) > 100000 {
		return NoteContent(""), ErrContentTooLong
	}

	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(normalizedValue)

	return NoteContent(sanitized), nil
}
