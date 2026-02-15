package application

import (
	"context"

	"pinnado/internal/shared/application/dto"
	"pinnado/internal/shared/domain"
)

type healthService struct {
	mongoClient MongoClient
}

func NewHealthService(mongoClient MongoClient) *healthService {
	return &healthService{
		mongoClient: mongoClient,
	}
}

func (s *healthService) HealthCheck(ctx context.Context) dto.HealthOutput {
	if s.mongoClient == nil {
		healthStatus := domain.HealthStatusError()
		return dto.HealthOutput{
			Status: healthStatus.Status,
		}
	}

	if err := s.mongoClient.Ping(ctx, nil); err != nil {
		healthStatus := domain.HealthStatusError()
		return dto.HealthOutput{
			Status: healthStatus.Status,
		}
	}

	healthStatus := domain.HealthStatusOk()
	return dto.HealthOutput{
		Status: healthStatus.Status,
	}
}
