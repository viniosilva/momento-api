package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"momento/internal/shared/app"
	"momento/internal/shared/domain"
)

type mockPinger struct {
	err error
}

func (m mockPinger) PingContext(_ context.Context) error {
	return m.err
}

func TestHealthService_HealthCheck(t *testing.T) {
	t.Run("should return ok status when db ping succeeds", func(t *testing.T) {
		pinger := mockPinger{}
		healthService := app.NewHealthService(pinger)

		got := healthService.HealthCheck(t.Context())

		assert.Equal(t, domain.HealthStatusValueOk, got.Status)
	})

	t.Run("should return error status when db client is nil", func(t *testing.T) {
		healthService := app.NewHealthService(nil)

		got := healthService.HealthCheck(t.Context())

		assert.Equal(t, domain.HealthStatusValueError, got.Status)
	})

	t.Run("should return error status when db ping fails", func(t *testing.T) {
		pinger := mockPinger{err: errors.New("connection failed")}
		healthService := app.NewHealthService(pinger)

		got := healthService.HealthCheck(t.Context())

		assert.Equal(t, domain.HealthStatusValueError, got.Status)
	})
}
