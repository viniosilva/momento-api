package presentation_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authinfra "pinnado/internal/auth/infrastructure"
	"pinnado/internal/notes/presentation"
	"pinnado/mocks"
	"pinnado/pkg/logger"

	"github.com/stretchr/testify/assert"
)

const (
	secretTest     = "secretTest"
	expirationTest = 5 * time.Minute
)

func TestSetupRouter(t *testing.T) {
	t.Run("should register POST /api/notes route", func(t *testing.T) {
		mux := http.NewServeMux()
		mockService := mocks.NewMockNoteService(t)
		jwtService := authinfra.NewJWTService(secretTest, expirationTest)

		presentation.SetupRouter(presentation.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			NoteService: mockService,
			JWTService:  jwtService,
			Logger:      logger.NewLogger("info"),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/notes", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		mux := http.NewServeMux()
		mockService := mocks.NewMockNoteService(t)
		jwtService := authinfra.NewJWTService(secretTest, expirationTest)

		assert.NotPanics(t, func() {
			presentation.SetupRouter(presentation.SetupRouterOptions{
				Mux:         mux,
				Prefix:      "/api",
				NoteService: mockService,
				JWTService:  jwtService,
				Logger:      nil,
			})
		})
	})

	t.Run("should apply logging middleware when logger is provided", func(t *testing.T) {
		mux := http.NewServeMux()
		mockService := mocks.NewMockNoteService(t)
		jwtService := authinfra.NewJWTService(secretTest, expirationTest)
		appLogger := slog.Default()

		presentation.SetupRouter(presentation.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			NoteService: mockService,
			JWTService:  jwtService,
			Logger:      appLogger,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/notes", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})
}
