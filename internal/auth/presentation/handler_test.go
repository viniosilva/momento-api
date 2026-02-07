package presentation_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/application"
	"pinnado/internal/auth/domain"
	"pinnado/internal/auth/infrastructure"
	"pinnado/internal/auth/presentation"
	"pinnado/mocks"
	"pinnado/pkg/nethttp"
)

// mapErrorToHTTPStatus is exported for testing
var mapErrorToHTTPStatus = presentation.MapErrorToHTTPStatus

const (
	secretTest     = "secretTest"
	expirationTest = 5 * time.Second
)

func TestNewAuthHandler(t *testing.T) {
	t.Run("should create auth handler", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
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
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
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
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		resp, got, err := nethttp.RequestWithResponse[string, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", "invalid json", handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
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
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
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
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(true, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Equal(t, "user already exists", got.Message)
	})

	t.Run("should return internal server error when repository fails", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
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
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)

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
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)

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

	t.Run("should return unauthorized when credentials are invalid", func(t *testing.T) {
		statusCode, message := mapErrorToHTTPStatus(domain.ErrInvalidCredentials)

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid credentials", message)
	})

	t.Run("should return internal server error when unknown fails", func(t *testing.T) {
		statusCode, message := mapErrorToHTTPStatus(context.DeadlineExceeded)

		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "internal server error", message)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	defaultReqBody := map[string]any{
		"email":    "user@example.com",
		"password": "ValidPass123!",
	}

	t.Run("should return ok when login is successful", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		password, err := domain.NewPassword(defaultReqBody["password"].(string))
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(user, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.LoginResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, got.Token)
	})

	t.Run("should return bad request when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		resp, got, err := nethttp.RequestWithResponse[string, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", "invalid json", handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		reqBody := map[string]any{
			"email":    "invalid-email",
			"password": defaultReqBody["password"],
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", reqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "invalid email format")
	})

	t.Run("should return unauthorized when user is not found", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, domain.ErrUserNotFound).
			Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "invalid credentials", got.Message)
	})

	t.Run("should return unauthorized when password is incorrect", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		password, err := domain.NewPassword("OtherPass123!")
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(user, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "invalid credentials", got.Message)
	})

	t.Run("should return internal server error when repository fails", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := infrastructure.NewJWTService(secretTest, expirationTest)
		authService := application.NewAuthService(mockRepo, jwtService)
		handler := presentation.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, assert.AnError).
			Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}
