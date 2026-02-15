package dto

import "pinnado/internal/shared/domain"

type HealthOutput struct {
	Status domain.HealthStatusEnum
}
