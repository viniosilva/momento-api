package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"momento/pkg/postgres"
)

func TestConnect(t *testing.T) {
	dsn := "postgres://user:pass@localhost:5433/nonexistent"
	maxRetries := 2
	retryDelay := 1 * time.Millisecond
	connectTimeout := 1 * time.Second

	t.Run("should retry on connection failure", func(t *testing.T) {
		got, err := postgres.Connect(
			t.Context(),
			dsn,
			maxRetries,
			retryDelay,
			connectTimeout,
		)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "failed to connect to PostgreSQL after 2 attempts")
	})
}
