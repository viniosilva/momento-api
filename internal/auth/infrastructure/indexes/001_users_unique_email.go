package indexes

import (
	"context"
	"momento/internal/auth/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUserEmailIndex(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection(domain.UsersCollectionName)

	indexModel := mongo.IndexModel{
		Keys:    map[string]any{"email": 1},
		Options: options.Index().SetUnique(true).SetName("unique_email"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
