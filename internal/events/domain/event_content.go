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

type EventContent string

func NewEventContent(value string) (EventContent, error) {
	normalizedValue := strings.TrimSpace(value)

	if normalizedValue == "" {
		return EventContent(""), ErrContentEmpty
	}

	if len(normalizedValue) > 100000 {
		return EventContent(""), ErrContentTooLong
	}

	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(normalizedValue)

	return EventContent(sanitized), nil
}