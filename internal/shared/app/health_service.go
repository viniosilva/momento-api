package app

import (
	"context"

	"momento/internal/shared/domain"
)

type healthService struct {
	mongoClient MongoClient
}

func NewHealthService(mongoClient MongoClient) *healthService {
	return &healthService{
		mongoClient: mongoClient,
	}
}

func (s *healthService) HealthCheck(ctx context.Context) HealthOutput {
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
