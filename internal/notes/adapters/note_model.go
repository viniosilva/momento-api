package adapters

import (
	"fmt"
	"time"

	"momento/internal/notes/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const notesCollectionName = "notes"

type noteDocument struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	ArchivedAt *time.Time         `bson:"archived_at"`
}

func toNoteDocument(n domain.Note) (noteDocument, error) {
	id, err := primitive.ObjectIDFromHex(n.ID)
	if err != nil {
		return noteDocument{}, fmt.Errorf("invalid note ID: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(n.UserID)
	if err != nil {
		return noteDocument{}, fmt.Errorf("invalid user ID: %w", err)
	}

	return noteDocument{
		ID:         id,
		UserID:     userID,
		Title:      string(n.Title),
		Content:    string(n.Content),
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
		ArchivedAt: n.ArchivedAt,
	}, nil
}

func toNoteDomain(d noteDocument) domain.Note {
	return domain.Note{
		ID:         d.ID.Hex(),
		UserID:     d.UserID.Hex(),
		Title:      domain.NoteTitle(d.Title),
		Content:    domain.NoteContent(d.Content),
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
		ArchivedAt: d.ArchivedAt,
	}
}

func parseObjectID(hex string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid ID %q: %w", hex, err)
	}
	return id, nil
}
