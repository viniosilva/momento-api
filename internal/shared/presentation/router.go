package presentation

import (
	"fmt"
	"net/http"

	"pinnado/internal/shared/application"
)

func SetupHealthRouter(mux *http.ServeMux, prefix string, healthService *application.HealthService) {
	handler := NewHealthHandler(healthService)

	mux.HandleFunc(fmt.Sprintf("GET %s/health", prefix), handler.HealthCheck)
}
