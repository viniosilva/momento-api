package presentation

import (
	"context"
	shareddto "momento/internal/shared/application/dto"
)

type HealthService interface {
	HealthCheck(ctx context.Context) shareddto.HealthOutput
}
