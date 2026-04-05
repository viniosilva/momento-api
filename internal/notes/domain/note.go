package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const NotesCollectionName = "notes"

var (
	ErrNoteNotFound       = errors.New("note not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Note struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Title      NoteTitle          `bson:"title"`
	Content    NoteContent        `bson:"content"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	ArchivedAt *time.Time         `bson:"archived_at"`
}

func NewNote(userID primitive.ObjectID, title NoteTitle, content NoteContent) Note {
	now := time.Now().UTC()

	return Note{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (n *Note) SetTitle(title NoteTitle) {
	n.Title = title
}

func (n *Note) SetContent(content NoteContent) {
	n.Content = content
}
