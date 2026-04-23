package indexes

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateEventUserIDIndex(ctx context.Context, db *mongo.Database, collectionName string) error {
	collection := db.Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    map[string]any{"user_id": 1},
		Options: options.Index().SetName("idx_user_id"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}