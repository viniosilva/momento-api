package application

import (
	"context"

	"pinnado/internal/shared/domain"
)

type HealthService struct {
	mongoClient MongoClient
}

func NewHealthService(mongoClient MongoClient) *HealthService {
	return &HealthService{
		mongoClient: mongoClient,
	}
}

func (s *HealthService) HealthCheck(ctx context.Context) HealthOutput {
	if s.mongoClient == nil {
		healthStatus := domain.HealthStatusError()
		return HealthOutput{
			Status: healthStatus.Status,
		}
	}

	if err := s.mongoClient.Ping(ctx, nil); err != nil {
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
