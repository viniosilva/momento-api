package presentation_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/presentation"
	"pinnado/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("should serve swagger.json file", func(t *testing.T) {
		// Create a temporary swagger.json file for testing
		tmpDir := t.TempDir()
		docsDir := filepath.Join(tmpDir, "docs")
		err := os.MkdirAll(docsDir, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(docsDir)

		swaggerJSONPath := filepath.Join(docsDir, "swagger.json")

		swaggerContent := map[string]any{
			"swagger": "2.0",
			"info": map[string]any{
				"title":   "Test API",
				"version": "1.0",
			},
			"paths": map[string]any{},
		}

		content, err := json.Marshal(swaggerContent)
		require.NoError(t, err)

		err = os.WriteFile(swaggerJSONPath, content, 0644)
		require.NoError(t, err)

		// Change to temp directory so relative path works
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalDir)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		mux := http.NewServeMux()
		healthService := application.NewHealthService(nil)
		appLogger := logger.NewLogger("info")

		presentation.SetupRouter(presentation.SetupRouterOptions{
			Mux:           mux,
			Prefix:        "/api",
			HealthService: healthService,
			Logger:        appLogger,
		})

		req := httptest.NewRequest(http.MethodGet, "/docs/swagger.json", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		var swaggerDoc map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &swaggerDoc)
		require.NoError(t, err, "swagger.json should be valid JSON")

		assert.Equal(t, "2.0", swaggerDoc["swagger"])
		assert.NotNil(t, swaggerDoc["info"])
		assert.NotNil(t, swaggerDoc["paths"])
	})
}
