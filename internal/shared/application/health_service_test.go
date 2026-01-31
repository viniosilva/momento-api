package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"pinnado/internal/shared/application"
)

func TestHealthService_HealthCheck(t *testing.T) {
	healthService := application.NewHealthService()

	t.Run("should return health output with ok status", func(t *testing.T) {
		got := healthService.HealthCheck(context.Background())

		assert.Equal(t, "ok", got.Status)
	})
}
