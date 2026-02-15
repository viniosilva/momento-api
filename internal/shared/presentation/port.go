package presentation

import (
	"context"
	shareddto "pinnado/internal/shared/application/dto"
)

type HealthService interface {
	HealthCheck(ctx context.Context) shareddto.HealthOutput
}
