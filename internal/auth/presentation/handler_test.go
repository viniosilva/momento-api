package presentation_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/application"
	"pinnado/internal/auth/domain"
	"pinnado/internal/auth/presentation"
	"pinnado/mocks"
	"pinnado/pkg/nethttp"
)

// mapErrorToHTTPStatus is exported for testing
var mapErrorToHTTPStatus = presentation.MapErrorToHTTPStatus

func TestNewAuthHandler(t *testing.T) {
	t.Run("should create auth handler", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		assert.NotNil(t, handler)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	defaultReqBody := map[string]any{
		"email":    "user@example.com",
		"password": "ValidPass123!",
	}

	t.Run("should return created when user is created successfully", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, nil).Once()
		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.RegisterResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, defaultReqBody["email"], got.Email)
	})

	t.Run("should return bad request when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		resp, got, err := nethttp.RequestWithResponse[string, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", "invalid json", handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		reqBody := map[string]any{
			"email":    "   ",
			"password": defaultReqBody["password"],
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", reqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "email cannot be empty")
	})

	t.Run("should return bad request when password is empty after trim", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		reqBody := map[string]any{
			"email":    defaultReqBody["email"],
			"password": "   ",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", reqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "password must be at least 8 characters", got.Message)
	})

	t.Run("should return conflict when user already exists", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(true, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Equal(t, "user already exists", got.Message)
	})

	t.Run("should return internal server error when repository error occurs", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)
		handler := presentation.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, assert.AnError).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestSetupRouter(t *testing.T) {
	t.Run("should setup router with auth handler", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)

		mux := http.NewServeMux()
		presentation.SetupRouter(presentation.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      nil,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		// Should not return 404 (route exists)
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	t.Run("should apply logging middleware when logger is provided", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(mockRepo)

		mux := http.NewServeMux()
		logger := slog.Default()
		presentation.SetupRouter(presentation.SetupRouterOptions{
			Mux:         mux,
			Prefix:      "/api",
			AuthService: authService,
			Logger:      logger,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader([]byte("{}")))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		// Should not return 404 (route exists)
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

func TestMapErrorToHTTPStatus(t *testing.T) {
	t.Run("should return conflict when user already exists", func(t *testing.T) {
		statusCode, message := mapErrorToHTTPStatus(domain.ErrUserAlreadyExists)

		assert.Equal(t, http.StatusConflict, statusCode)
		assert.Equal(t, "user already exists", message)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		testCases := []struct {
			err      error
			expected string
		}{
			{domain.ErrEmailIsEmpty, "email cannot be empty"},
			{domain.ErrInvalidEmail, "invalid email format"},
		}

		for _, tc := range testCases {
			statusCode, message := mapErrorToHTTPStatus(tc.err)

			assert.Equal(t, http.StatusBadRequest, statusCode)
			assert.Equal(t, tc.expected, message)
		}
	})

	t.Run("should return bad request when password is invalid", func(t *testing.T) {
		testCases := []struct {
			err      error
			expected string
		}{
			{domain.ErrPasswordTooShort, "password must be at least 8 characters"},
			{domain.ErrPasswordTooLong, "password must be less than 64 characters"},
			{domain.ErrPasswordMissingUpper, "password must contain at least one uppercase letter"},
			{domain.ErrPasswordMissingLower, "password must contain at least one lowercase letter"},
			{domain.ErrPasswordMissingNumber, "password must contain at least one number"},
			{domain.ErrPasswordMissingSymbol, "password must contain at least one symbol"},
		}

		for _, tc := range testCases {
			statusCode, message := mapErrorToHTTPStatus(tc.err)

			assert.Equal(t, http.StatusBadRequest, statusCode)
			assert.Equal(t, tc.expected, message)
		}
	})

	t.Run("should return internal server error when unknown error occurs", func(t *testing.T) {
		statusCode, message := mapErrorToHTTPStatus(context.DeadlineExceeded)

		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "internal server error", message)
	})
}
