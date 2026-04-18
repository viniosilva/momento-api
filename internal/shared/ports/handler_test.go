package ports_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"momento/internal/shared/app"
	"momento/internal/shared/domain"
	"momento/internal/shared/mocks"
	"momento/internal/shared/ports"
	"momento/internal/shared/ports/response"
	"momento/pkg/nethttp"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	t.Run("should return healthcheck ok response", func(t *testing.T) {
		mockMongoClient := mocks.NewMockMongoClient(t)
		mockMongoClient.EXPECT().Ping(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		healthService := app.NewHealthService(mockMongoClient)
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
		mockMongoClient := mocks.NewMockMongoClient(t)
		mockMongoClient.EXPECT().Ping(mock.Anything, mock.Anything).
			Return(errors.New("connection failed")).
			Once()

		healthService := app.NewHealthService(mockMongoClient)
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
