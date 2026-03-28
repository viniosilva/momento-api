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
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Content   NoteContent        `bson:"content"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewNote(userID primitive.ObjectID, content NoteContent) Note {
	now := time.Now().UTC()

	return Note{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
