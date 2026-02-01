package presentation_test

import (
	"net/http"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/presentation"
	"pinnado/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	t.Run("should not panic when setting up health router", func(t *testing.T) {
		mux := http.NewServeMux()
		healthService := application.NewHealthService(nil)
		prefix := "/api"
		appLogger := logger.NewLogger("info")

		assert.NotPanics(t, func() {
			presentation.SetupRouter(presentation.SetupRouterOptions{
				Mux:           mux,
				Prefix:        prefix,
				HealthService: healthService,
				Logger:        appLogger,
			})
		})
	})

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		mux := http.NewServeMux()
		healthService := application.NewHealthService(nil)
		prefix := "/api"

		assert.NotPanics(t, func() {
			presentation.SetupRouter(presentation.SetupRouterOptions{
				Mux:           mux,
				Prefix:        prefix,
				HealthService: healthService,
				Logger:        nil,
			})
		})
	})
}
