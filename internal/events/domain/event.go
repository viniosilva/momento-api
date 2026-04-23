package domain

import (
	"errors"
	"time"

	"momento/pkg/uid"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	ID          string
	OwnerUserID string
	Title       EventTitle
	Content     EventContent
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

func NewEvent(ownerUserID string, title EventTitle, content EventContent) Event {
	now := time.Now().UTC()

	return Event{
		ID:          uid.New(),
		OwnerUserID: ownerUserID,
		Title:       title,
		Content:     content,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (e *Event) Update(title EventTitle, content EventContent) {
	e.Title = title
	e.Content = content
	e.UpdatedAt = time.Now().UTC()
}
