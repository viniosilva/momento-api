package application

import (
	"pinnado/internal/notes/domain"
	"time"
)

type NoteInput struct {
	UserID  string
	Content string
}

type NoteOutput struct {
	ID        string
	UserID    string
	Content   domain.NoteContent
	CreatedAt time.Time
	UpdatedAt time.Time
}
