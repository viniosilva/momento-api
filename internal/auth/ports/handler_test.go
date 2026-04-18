package ports_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"momento/internal/auth/adapters"
	"momento/internal/auth/app"
	"momento/internal/auth/domain"
	"momento/internal/auth/mocks"
	"momento/internal/auth/ports"
	"momento/pkg/nethttp"
)

const (
	secretTest     = "secretTest"
	expirationTest = 5 * time.Second
)

func TestNewAuthHandler(t *testing.T) {
	t.Run("should create auth handler", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

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
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, nil).Once()
		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, ports.RegisterResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, defaultReqBody["email"], got.Email)
	})

	t.Run("should return bad request when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		resp, got, err := nethttp.RequestWithResponse[string, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", "invalid json", handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		reqBody := map[string]any{
			"email":    "   ",
			"password": defaultReqBody["password"],
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", reqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "email cannot be empty")
	})

	t.Run("should return bad request when password is empty after trim", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		reqBody := map[string]any{
			"email":    defaultReqBody["email"],
			"password": "   ",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", reqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "password must be at least 6 characters", got.Message)
	})

	t.Run("should return conflict when user already exists", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(true, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Equal(t, "user already exists", got.Message)
	})

	t.Run("should return internal server error when repository fails", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		mockRepo.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, assert.AnError).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register", defaultReqBody, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestMapErrorToHTTPStatus(t *testing.T) {
	t.Run("should return conflict when user already exists", func(t *testing.T) {
		statusCode, message := ports.MapErrorToHTTPStatus(domain.ErrUserAlreadyExists)

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
			statusCode, message := ports.MapErrorToHTTPStatus(tc.err)

			assert.Equal(t, http.StatusBadRequest, statusCode)
			assert.Equal(t, tc.expected, message)
		}
	})

	t.Run("should return bad request when password is invalid", func(t *testing.T) {
		testCases := []struct {
			err      error
			expected string
		}{
			{domain.ErrPasswordTooShort, "password must be at least 6 characters"},
			{domain.ErrPasswordTooLong, "password must be less than 64 characters"},
			{domain.ErrPasswordMissingUpper, "password must contain at least one uppercase letter"},
			{domain.ErrPasswordMissingLower, "password must contain at least one lowercase letter"},
			{domain.ErrPasswordMissingNumber, "password must contain at least one number"},
			{domain.ErrPasswordMissingSymbol, "password must contain at least one symbol"},
		}

		for _, tc := range testCases {
			statusCode, message := ports.MapErrorToHTTPStatus(tc.err)

			assert.Equal(t, http.StatusBadRequest, statusCode)
			assert.Equal(t, tc.expected, message)
		}
	})

	t.Run("should return unauthorized when credentials are invalid", func(t *testing.T) {
		statusCode, message := ports.MapErrorToHTTPStatus(domain.ErrInvalidCredentials)

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid credentials", message)
	})

	t.Run("should return internal server error when unknown fails", func(t *testing.T) {
		statusCode, message := ports.MapErrorToHTTPStatus(context.DeadlineExceeded)

		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "internal server error", message)
	})

	t.Run("should not leak internal wrapper text on unknown error", func(t *testing.T) {
		wrapped := fmt.Errorf("s.userRepository.FindByEmail: %w", assert.AnError)
		code, msg := ports.MapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, "internal server error", msg)
		assert.NotContains(t, msg, "userRepository")
		assert.NotContains(t, msg, "FindByEmail")
	})

	t.Run("should use canonical domain message when sentinel is wrapped", func(t *testing.T) {
		wrapped := fmt.Errorf("app.Register: %w", domain.ErrInvalidEmail)
		code, msg := ports.MapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, domain.ErrInvalidEmail.Error(), msg)
		assert.NotContains(t, msg, "app.Register")
	})
}

func TestAuthHandler_Login(t *testing.T) {
	defaultReqBody := map[string]any{
		"email":    "user@example.com",
		"password": "ValidPass123!",
	}

	t.Run("should return ok when login is successful", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		password, err := domain.NewPassword(defaultReqBody["password"].(string))
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(user, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, ports.LoginResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, got.Token)
	})

	t.Run("should return bad request when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		resp, got, err := nethttp.RequestWithResponse[string, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", "invalid json", handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		reqBody := map[string]any{
			"email":    "invalid-email",
			"password": defaultReqBody["password"],
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", reqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "invalid email format")
	})

	t.Run("should return unauthorized when user is not found", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, domain.ErrUserNotFound).
			Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "invalid credentials", got.Message)
	})

	t.Run("should return unauthorized when password is incorrect", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		password, err := domain.NewPassword("OtherPass123!")
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(user, nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "invalid credentials", got.Message)
	})

	t.Run("should return internal server error when repository fails", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		authService := app.NewAuthService(mockRepo, jwtService)
		handler := ports.NewAuthHandler(authService)

		email, err := domain.NewEmail(defaultReqBody["email"].(string))
		require.NoError(t, err)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, assert.AnError).
			Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login", defaultReqBody, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}
