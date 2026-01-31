package presentation_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"pinnado/internal/shared/application"
	"pinnado/internal/shared/presentation"
)

func TestSetupHealthRouter(t *testing.T) {
	t.Run("should not panic when setting up health router", func(t *testing.T) {
		mux := http.NewServeMux()
		healthService := application.NewHealthService()
		prefix := "/api"

		assert.NotPanics(t, func() {
			presentation.SetupHealthRouter(mux, prefix, healthService)
		})
	})
}
