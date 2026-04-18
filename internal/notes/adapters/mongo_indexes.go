package adapters

import (
	"context"
	"fmt"

	"momento/internal/notes/adapters/indexes"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	if db == nil {
		return fmt.Errorf("mongo client is nil")
	}

	if err := indexes.CreateNoteUserIDIndex(ctx, db); err != nil {
		return fmt.Errorf("failed to create index on user_id field: %w", err)
	}

	return nil
}
