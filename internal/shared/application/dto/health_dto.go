package dto

import "momento/internal/shared/domain"

type HealthOutput struct {
	Status domain.HealthStatusEnum
}
