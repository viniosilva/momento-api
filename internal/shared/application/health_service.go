package application

import (
	"context"

	"pinnado/internal/shared/domain"
)

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) HealthCheck(ctx context.Context) HealthOutput {
	healthStatus := domain.HealthStatusOk()

	return HealthOutput{
		Status: string(healthStatus.Status),
	}
}
