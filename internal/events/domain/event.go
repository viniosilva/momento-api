package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	ID          string
	OwnerUserID string
	Title       EventTitle
	Content     EventContent
	Metadata    *EventMetadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

func NewEvent(ownerUserID string, title EventTitle, content EventContent) Event {
	now := time.Now().UTC()
	metadata := NewEventMetadata()

	return Event{
		ID:          uuid.NewString(),
		OwnerUserID: ownerUserID,
		Title:       title,
		Content:     content,
		Metadata:    &metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (e *Event) Update(title EventTitle, content EventContent) {
	e.Title = title
	e.Content = content
	e.UpdatedAt = time.Now().UTC()
}

func (e *Event) AddImage(path ImagePath) error {
	if e.Metadata == nil {
		e.Metadata = new(NewEventMetadata())
	}

	e.UpdatedAt = time.Now().UTC()

	return e.Metadata.AddImage(path)
}

func (e *Event) RemoveImage(path ImagePath) error {
	if e.Metadata == nil {
		return ErrImageNotFound
	}

	e.UpdatedAt = time.Now().UTC()

	return e.Metadata.RemoveImage(path)
}
