package domain

import (
	"errors"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	ErrTitleEmpty   = errors.New("title cannot be empty")
	ErrTitleTooLong = errors.New("title exceeds maximum length of 24 characters")
)

type EventTitle string

func NewEventTitle(value string) (EventTitle, error) {
	normalizedValue := strings.TrimSpace(value)

	if normalizedValue == "" {
		return EventTitle(""), ErrTitleEmpty
	}

	if len(normalizedValue) > 24 {
		return EventTitle(""), ErrTitleTooLong
	}

	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(normalizedValue)

	return EventTitle(sanitized), nil
}