package ports_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authadapters "momento/internal/auth/adapters"
	"momento/internal/events/mocks"
	"momento/internal/events/ports"
	"momento/pkg/logger"

	"github.com/stretchr/testify/assert"
)

const (
	secretTest     = "secretTest"
	expirationTest = 5 * time.Minute
)

func TestSetupRouter(t *testing.T) {
	t.Run("should register POST /api/events route", func(t *testing.T) {
		mux := http.NewServeMux()
		mockService := mocks.NewMockEventService(t)
		jwtService := authadapters.NewJWTService(secretTest, expirationTest)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:          mux,
			Prefix:       "/api",
			EventService: mockService,
			JWTService:   jwtService,
			Logger:       logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/events", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		mux := http.NewServeMux()
		mockService := mocks.NewMockEventService(t)
		jwtService := authadapters.NewJWTService(secretTest, expirationTest)

		assert.NotPanics(t, func() {
			ports.SetupRouter(ports.SetupRouterOptions{
				Mux:          mux,
				Prefix:       "/api",
				EventService: mockService,
				JWTService:   jwtService,
				Logger:       nil,
			})
		})
	})

	t.Run("should apply logging middleware when logger is provided", func(t *testing.T) {
		mux := http.NewServeMux()
		mockService := mocks.NewMockEventService(t)
		jwtService := authadapters.NewJWTService(secretTest, expirationTest)
		appLogger := slog.Default()

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:          mux,
			Prefix:       "/api",
			EventService: mockService,
			JWTService:   jwtService,
			Logger:       appLogger,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/events", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})
}
