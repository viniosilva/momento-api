package presentation

import (
	"net/http"
	"momento/internal/shared/domain"
	"momento/internal/shared/presentation/response"
	nethttp_utils "momento/pkg/nethttp/utils"
)

type healthHandler struct {
	healthService HealthService
}

func NewHealthHandler(healthService HealthService) *healthHandler {
	return &healthHandler{
		healthService: healthService,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application and database connection
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.HealthResponse "Health status response"
// @Failure 503 {object} response.HealthResponse "Health status response"
// @Router /api/healthcheck [get]
func (h *healthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	output := h.healthService.HealthCheck(r.Context())

	response := response.HealthResponse{
		Status: output.Status,
	}
	statusCode := http.StatusOK
	if output.Status == domain.HealthStatusValueError {
		statusCode = http.StatusServiceUnavailable
	}

	nethttp_utils.JSON(w, statusCode, response)
}
