package ports_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"momento/internal/auth/adapters"
	"momento/internal/auth/app"
	"momento/internal/auth/mocks"
	"momento/internal/auth/ports"
	"momento/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newAuthService(t *testing.T) ports.AuthService {
	t.Helper()
	userRepo := mocks.NewMockUserRepository(t)
	secureTokenSvc := mocks.NewMockSecureTokenService(t)
	jwtService := adapters.NewJWTService(secretTest, expirationTest)
	resetTokenSvc := mocks.NewMockResetTokenService(t)
	emailSender := mocks.NewMockEmailSender(t)

	resetTokenSize := 32
	resetTokenTTL := 1 * time.Hour

	return app.NewAuthService(userRepo, jwtService, secureTokenSvc, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)
}

func TestSetupRouter(t *testing.T) {
	t.Run("should register POST /api/auth/register route", func(t *testing.T) {
		mux := http.NewServeMux()
		authService := newAuthService(t)

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
		authService := newAuthService(t)

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

	t.Run("should register POST /api/auth/refresh route", func(t *testing.T) {
		mux := http.NewServeMux()
		authService := newAuthService(t)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should register POST /api/auth/logout route", func(t *testing.T) {
		mux := http.NewServeMux()
		authService := newAuthService(t)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should register POST /api/auth/forgot-password route", func(t *testing.T) {
		mux := http.NewServeMux()
		authService := newAuthService(t)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should register POST /api/auth/reset-password route", func(t *testing.T) {
		mux := http.NewServeMux()
		authService := newAuthService(t)

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/reset-password", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should register GET /api/auth/reset-password/validate route", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)

		resetTokenSize := 32
		resetTokenTTL := 1 * time.Hour

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		mux := http.NewServeMux()

		ports.SetupRouter(ports.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodGet, "/api/auth/reset-password/validate?token=test", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		mux := http.NewServeMux()
		authService := newAuthService(t)

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
		authService := newAuthService(t)
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
