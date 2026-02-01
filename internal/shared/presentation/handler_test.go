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

		requestBody := map[string]any{}
		cb := func(w http.ResponseWriter, r *http.Request) error {
			handler.HealthCheck(w, r)
			return nil
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]interface{}, presentation.HealthResponse](
			t.Context(),
			http.MethodGet,
			"/api/healthcheck",
			requestBody,
			cb,
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

		requestBody := map[string]any{}
		cb := func(w http.ResponseWriter, r *http.Request) error {
			handler.HealthCheck(w, r)
			return nil
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]interface{}, presentation.HealthResponse](
			t.Context(),
			http.MethodGet,
			"/api/healthcheck",
			requestBody,
			cb,
		)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, domain.HealthStatusValueError, got.Status)
	})
}
