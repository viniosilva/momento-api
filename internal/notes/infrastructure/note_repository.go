package infrastructure

import (
	"context"

	"pinnado/internal/notes/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type noteRepository struct {
	collection *mongo.Collection
}

func NewNoteRepository(collection *mongo.Collection) *noteRepository {
	return &noteRepository{
		collection: collection,
	}
}

func (r *noteRepository) Create(ctx context.Context, note domain.Note) error {
	_, err := r.collection.InsertOne(ctx, note)
	return err
}
