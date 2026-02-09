package nethttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/domain"
	"pinnado/internal/auth/infrastructure"
	"pinnado/pkg/nethttp"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key"
	expiration := 5 * time.Minute
	jwtService := infrastructure.NewJWTService(secret, expiration)

	t.Run("should allow request with valid token", func(t *testing.T) {
		middleware := nethttp.AuthMiddleware(jwtService)

		userID := "507f1f77bcf86cd799439011"
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		token, err := jwtService.Generate(userID, email)
		require.NoError(t, err)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxUserID := r.Context().Value(nethttp.ContextKeyUserID)
			ctxEmail := r.Context().Value(nethttp.ContextKeyEmail)

			assert.Equal(t, userID, ctxUserID)
			assert.Equal(t, email, ctxEmail)

			nethttp.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return 401 when authorization header is missing", func(t *testing.T) {
		middleware := nethttp.AuthMiddleware(jwtService)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		})

		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "authorization header is required")
	})

	t.Run("should return 401 when authorization header format is invalid", func(t *testing.T) {
		middleware := nethttp.AuthMiddleware(jwtService)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		})

		wrappedHandler := middleware(handler)

		testCases := []struct {
			name        string
			header      string
			expectedMsg string
		}{
			{"missing Bearer prefix", "token-value", "invalid authorization header format"},
			{"missing token", "Bearer ", "invalid authorization header format"},
			{"too many parts", "Bearer token extra", "invalid authorization header format"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.Header.Set("Authorization", tc.header)
				rec := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(rec, req)

				assert.Equal(t, http.StatusUnauthorized, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedMsg)
			})
		}
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		middleware := nethttp.AuthMiddleware(jwtService)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		})

		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid or expired token")
	})

	t.Run("should return 401 when token is expired", func(t *testing.T) {
		expiredJWTService := infrastructure.NewJWTService(secret, -5*time.Second)
		middleware := nethttp.AuthMiddleware(expiredJWTService)

		userID := "user123"
		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		expiredToken, err := expiredJWTService.Generate(userID, email)
		require.NoError(t, err)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		})

		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid or expired token")
	})
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	secret := "test-secret-key"
	expiration := 5 * time.Minute
	jwtService := infrastructure.NewJWTService(secret, expiration)

	t.Run("should inject UserID and Email into context", func(t *testing.T) {
		middleware := nethttp.AuthMiddleware(jwtService)

		userID := "user123"
		email, err := domain.NewEmail("test@example.com")
		require.NoError(t, err)

		token, err := jwtService.Generate(userID, email)
		require.NoError(t, err)

		var capturedUserID, capturedEmail interface{}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedUserID = r.Context().Value(nethttp.ContextKeyUserID)
			capturedEmail = r.Context().Value(nethttp.ContextKeyEmail)

			nethttp.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		assert.Equal(t, "user123", capturedUserID)
		assert.Equal(t, domain.Email("test@example.com"), capturedEmail)
	})
}
