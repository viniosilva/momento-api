package ports_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"momento/internal/shared/app"
	"momento/internal/shared/domain"
	"momento/internal/shared/ports"
	"momento/internal/shared/ports/response"
	"momento/pkg/nethttp"
)

type mockPinger struct {
	err error
}

func (m mockPinger) PingContext(_ context.Context) error {
	return m.err
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	t.Run("should return healthcheck ok response", func(t *testing.T) {
		pinger := mockPinger{}
		healthService := app.NewHealthService(pinger)
		handler := ports.NewHealthHandler(healthService)

		resp, got, err := nethttp.RequestWithResponse[map[string]any, response.HealthResponse](
			t.Context(),
			http.MethodGet,
			"/api/healthcheck",
			nil,
			handler.HealthCheck,
		)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, domain.HealthStatusValueOk, got.Status)
	})

	t.Run("should return healthcheck error response", func(t *testing.T) {
		pinger := mockPinger{err: errors.New("connection failed")}
		healthService := app.NewHealthService(pinger)
		handler := ports.NewHealthHandler(healthService)

		resp, got, err := nethttp.RequestWithResponse[map[string]any, response.HealthResponse](
			t.Context(),
			http.MethodGet,
			"/api/healthcheck",
			nil,
			handler.HealthCheck,
		)
		require.NoError(t, err)

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		assert.Equal(t, domain.HealthStatusValueError, got.Status)
	})
}
