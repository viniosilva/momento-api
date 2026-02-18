package nethttp_recovery_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	nethttp_recovery "pinnado/pkg/nethttp/recovery"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	logger := slog.Default()

	t.Run("should recover from panic and return internal server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		middleware := nethttp_recovery.RecoveryMiddleware(logger)
		recoveryHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		recoveryHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
	})

	t.Run("should recover from panic and return internal server error when logger is nil", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		middleware := nethttp_recovery.RecoveryMiddleware(nil)
		recoveryHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		recoveryHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
	})

	t.Run("should pass through when no panic occurs", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		middleware := nethttp_recovery.RecoveryMiddleware(logger)
		recoveryHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		recoveryHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})
}
