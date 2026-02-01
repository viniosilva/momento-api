package presentation

import (
	"net/http"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/domain"
	"pinnado/pkg/nethttp"
)

type HealthHandler struct {
	healthService *application.HealthService
}

func NewHealthHandler(healthService *application.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application and database connection
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "Health status response"
// @Failure 503 {object} HealthResponse "Health status response"
// @Router /healthcheck [get]
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	output := h.healthService.HealthCheck(r.Context())

	response := HealthResponse{
		Status: output.Status,
	}
	statusCode := http.StatusOK
	if output.Status == domain.HealthStatusValueError {
		statusCode = http.StatusServiceUnavailable
	}

	nethttp.JSON(w, statusCode, response)
}
