package presentation_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinnado/internal/shared/application"
	"pinnado/internal/shared/presentation"
	"pinnado/pkg/nethttp"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	healthService := application.NewHealthService()
	handler := presentation.NewHealthHandler(healthService)

	t.Run("should return healthcheck ok response", func(t *testing.T) {
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
		assert.Equal(t, "ok", got.Status)
	})
}
