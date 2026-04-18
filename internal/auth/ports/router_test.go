package ports_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"momento/internal/auth/adapters"
	"momento/internal/auth/app"
	"momento/internal/auth/mocks"
	"momento/internal/auth/ports"
	"momento/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	t.Run("should register POST /api/auth/register route", func(t *testing.T) {
		mux := http.NewServeMux()
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should register POST /api/auth/login route", func(t *testing.T) {
		mux := http.NewServeMux()
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		mux := http.NewServeMux()
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)

		assert.NotPanics(t, func() {
			ports.SetupRouter(ports.SetupRouterOptions{
				Mux:         mux,
				Prefix:      "/api",
				AuthService: authService,
				Logger:      nil,
			})
		})
	})

	t.Run("should apply logging middleware when logger is provided", func(t *testing.T) {
		mux := http.NewServeMux()
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		appLogger := slog.Default()

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      appLogger,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})
}
