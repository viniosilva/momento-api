package application_test

import (
	"errors"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/domain"
	"pinnado/internal/shared/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthService_HealthCheck(t *testing.T) {
	t.Run("should return ok status when mongo ping succeeds", func(t *testing.T) {
		mockMongoClient := mocks.NewMockMongoClient(t)
		mockMongoClient.EXPECT().Ping(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		healthService := application.NewHealthService(mockMongoClient)

		got := healthService.HealthCheck(t.Context())

		assert.Equal(t, domain.HealthStatusValueOk, got.Status)
	})

	t.Run("should return error status when mongo client is nil", func(t *testing.T) {
		healthService := application.NewHealthService(nil)

		got := healthService.HealthCheck(t.Context())

		assert.Equal(t, domain.HealthStatusValueError, got.Status)
	})

	t.Run("should return error status when mongo ping fails", func(t *testing.T) {
		mockMongoClient := mocks.NewMockMongoClient(t)
		mockMongoClient.EXPECT().Ping(mock.Anything, mock.Anything).
			Return(errors.New("connection failed")).
			Once()

		healthService := application.NewHealthService(mockMongoClient)

		got := healthService.HealthCheck(t.Context())

		assert.Equal(t, domain.HealthStatusValueError, got.Status)
	})
}
