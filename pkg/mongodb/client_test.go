package mongodb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"pinnado/pkg/mongodb"
)

func TestNewMongoClient(t *testing.T) {
	host := "localhost"
	port := "27000"
	dbName := "testdb"
	user := "user"
	pass := "pass"
	maxRetries := 2
	retryDelay := 1 * time.Millisecond
	connectTimeout := 1 * time.Second

	t.Run("should retry on connection failure", func(t *testing.T) {
		got, err := mongodb.NewMongoClient(
			t.Context(),
			host,
			port,
			dbName,
			user,
			pass,
			maxRetries,
			retryDelay,
			connectTimeout,
		)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "failed to connect to MongoDB after 2 attempts")
	})

	t.Run("should build connection URI without auth", func(t *testing.T) {
		got, err := mongodb.NewMongoClient(
			t.Context(),
			host,
			port,
			dbName,
			"",
			"",
			maxRetries,
			retryDelay,
			connectTimeout,
		)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "failed to connect to MongoDB after 2 attempts")
	})
}
