package presentation_test

import (
	"errors"
	"net/http"
	"momento/internal/shared/application"
	"momento/internal/shared/domain"
	"momento/internal/shared/mocks"
	"momento/internal/shared/presentation"
	"momento/internal/shared/presentation/response"
	"momento/pkg/nethttp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	t.Run("should return healthcheck ok response", func(t *testing.T) {
		mockMongoClient := mocks.NewMockMongoClient(t)
		mockMongoClient.EXPECT().Ping(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		healthService := application.NewHealthService(mockMongoClient)
		handler := presentation.NewHealthHandler(healthService)

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

		healthService := application.NewHealthService(mockMongoClient)
		handler := presentation.NewHealthHandler(healthService)

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
