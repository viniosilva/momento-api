package ports

import (
	"context"

	sharedapp "momento/internal/shared/app"
)

type HealthService interface {
	HealthCheck(ctx context.Context) sharedapp.HealthOutput
}
