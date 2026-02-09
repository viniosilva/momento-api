package indexes

import (
	"context"
	"pinnado/internal/notes/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateNoteUserIDIndex(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection(domain.NotesCollectionName)

	indexModel := mongo.IndexModel{
		Keys:    map[string]interface{}{"user_id": 1},
		Options: options.Index().SetName("idx_user_id"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
