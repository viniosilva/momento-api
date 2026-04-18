package response

import "momento/internal/shared/domain"

type HealthResponse struct {
	Status domain.HealthStatusEnum `json:"status" example:"ok"`
}
