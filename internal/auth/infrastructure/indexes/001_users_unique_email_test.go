package indexes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/infrastructure/indexes"
	"pinnado/pkg/mongodb"
)

// TestCreateUserEmailIndex tests the creation of unique email index.
// This is an integration test that requires a MongoDB instance.
// To run: go test -tags=integration ./internal/auth/infrastructure/indexes/...
func TestCreateUserEmailIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("should create unique index on email field", func(t *testing.T) {
		ctx := context.Background()
		mongoClient, err := mongodb.NewMongoClient(
			ctx,
			"localhost",
			"27017",
			"pinnado_test",
			"admin",
			"admin",
			3,
			2,
			10,
		)
		require.NoError(t, err)
		defer func() {
			if err := mongoClient.Disconnect(context.Background()); err != nil {
				t.Logf("error disconnecting from MongoDB: %v", err)
			}
		}()

		db := mongoClient.Database("pinnado_test")

		err = indexes.CreateUserEmailIndex(ctx, db)
		require.NoError(t, err)

		// Verify index was created by listing indexes
		collection := db.Collection("users")
		indexesList, err := collection.Indexes().List(ctx)
		require.NoError(t, err)

		var found bool
		for indexesList.Next(ctx) {
			var indexDoc map[string]interface{}
			err := indexesList.Decode(&indexDoc)
			require.NoError(t, err)

			if name, ok := indexDoc["name"].(string); ok && name == "unique_email" {
				found = true
				assert.True(t, indexDoc["unique"].(bool), "index should be unique")
				break
			}
		}

		assert.True(t, found, "unique_email index should be created")
	})
}
