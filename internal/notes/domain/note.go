package domain

import (
	"errors"
	"time"

	"momento/pkg/uid"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type Note struct {
	ID         string
	UserID     string
	Title      NoteTitle
	Content    NoteContent
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
}

func NewNote(userID string, title NoteTitle, content NoteContent) Note {
	now := time.Now().UTC()

	return Note{
		ID:        uid.New(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (n *Note) Update(title NoteTitle, content NoteContent) {
	n.Title = title
	n.Content = content
	n.UpdatedAt = time.Now().UTC()
}
