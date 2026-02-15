package presentation_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pinnado/internal/shared/application"
	"pinnado/internal/shared/domain"
	"pinnado/internal/shared/presentation"
	"pinnado/internal/shared/presentation/response"
	"pinnado/mocks"
	"pinnado/pkg/nethttp"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	t.Run("should return healthcheck ok response", func(t *testing.T) {
		mockMongoClient := mocks.NewMockMongoClient(t)
		mockMongoClient.EXPECT().Ping(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		healthService := application.NewHealthService(mockMongoClient)
		handler := presentation.NewHealthHandler(healthService)

		resp, got, err := nethttp.RequestWithResponse[map[string]interface{}, response.HealthResponse](
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

		healthService := application.NewHealthService(mockMongoClient)
		handler := presentation.NewHealthHandler(healthService)

		resp, got, err := nethttp.RequestWithResponse[map[string]interface{}, response.HealthResponse](
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
