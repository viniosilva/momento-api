package infrastructure

import (
	"context"
	"fmt"
	"pinnado/internal/auth/infrastructure/indexes"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	if db == nil {
		return fmt.Errorf("mongo client is nil")
	}

	if err := indexes.CreateUserEmailIndex(ctx, db); err != nil {
		return fmt.Errorf("failed to create unique index on email field: %w", err)
	}

	return nil
}
