package app

import (
	"context"

	"momento/internal/shared/domain"
)

type healthService struct {
	db Pinger
}

func NewHealthService(db Pinger) *healthService {
	return &healthService{
		db: db,
	}
}

func (s *healthService) HealthCheck(ctx context.Context) HealthOutput {
	if s.db == nil {
		healthStatus := domain.HealthStatusError()
		return HealthOutput{
			Status: healthStatus.Status,
		}
	}

	if err := s.db.PingContext(ctx); err != nil {
		healthStatus := domain.HealthStatusError()
		return HealthOutput{
			Status: healthStatus.Status,
		}
	}

	healthStatus := domain.HealthStatusOk()
	return HealthOutput{
		Status: healthStatus.Status,
	}
}
