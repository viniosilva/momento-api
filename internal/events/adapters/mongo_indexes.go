package adapters

import (
	"context"
	"fmt"

	"momento/internal/events/adapters/indexes"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	if db == nil {
		return fmt.Errorf("mongo client is nil")
	}

	if err := indexes.CreateEventUserIDIndex(ctx, db, eventsCollectionName); err != nil {
		return fmt.Errorf("failed to create index on user_id field: %w", err)
	}

	return nil
}