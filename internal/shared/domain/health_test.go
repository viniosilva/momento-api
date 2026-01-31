package domain_test

import (
	"pinnado/internal/shared/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthStatusOk(t *testing.T) {
	t.Run("should return health status with ok value", func(t *testing.T) {
		got := domain.HealthStatusOk()

		assert.Equal(t, domain.HealthStatusValueOk, got.Status)
	})
}
