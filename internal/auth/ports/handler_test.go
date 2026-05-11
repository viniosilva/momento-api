package ports_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	secretTest            = "secretTest"
	expirationTest        = 5 * time.Second
	resetTokenSize        = 32
	resetTokenTTL         = 1 * time.Hour
	verificationTokenSize = 32
	verificationTokenTTL  = 24 * time.Hour
	verificationURL       = "http://momentonow.com/verify-email"
)

func newAuthServiceWithMocks(t *testing.T) (ports.AuthService, *mocks.MockUserRepository, *mocks.MockSecureTokenService, *mocks.MockResetTokenService, *mocks.MockEmailSender) {
	t.Helper()
	userRepo := mocks.NewMockUserRepository(t)
	secureTokenSvc := mocks.NewMockSecureTokenService(t)
	jwtService := adapters.NewJWTService(secretTest, expirationTest)
	resetTokenSvc := mocks.NewMockResetTokenService(t)
	emailSender := mocks.NewMockEmailSender(t)
	tokenService := mocks.NewMockTokenService(t)

	svc := app.NewAuthService(userRepo, jwtService, secureTokenSvc, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
	return svc, userRepo, secureTokenSvc, resetTokenSvc, emailSender
}

func TestNewAuthHandler(t *testing.T) {
	t.Run("should create auth handler", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		assert.NotNil(t, handler)
	})
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	t.Run("should return no content when email is sent", func(t *testing.T) {
		svc, mockRepo, _, mockResetTokenSvc, mockEmailSender := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		email, _ := domain.NewEmail("user@example.com")
		password, _ := domain.NewPassword("ValidPass123!")
		user := domain.NewUser(email, password)

		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(user, nil).Once()
		mockResetTokenSvc.EXPECT().Store(mock.Anything, mock.Anything, user.ID, resetTokenTTL).Return(nil).Once()
		mockEmailSender.EXPECT().SendResetPasswordEmail(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := nethttp.Request(
			t.Context(), http.MethodPost, "/auth/forgot-password",
			ports.ForgotPasswordRequest{Email: "user@example.com"}, handler.ForgotPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should return no content even when user does not exist", func(t *testing.T) {
		svc, mockRepo, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		email, _ := domain.NewEmail("user@example.com")
		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(domain.User{}, domain.ErrUserNotFound).Once()

		resp, err := nethttp.Request(
			t.Context(), http.MethodPost, "/auth/forgot-password",
			ports.ForgotPasswordRequest{Email: "user@example.com"}, handler.ForgotPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[ports.ForgotPasswordRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/forgot-password",
			ports.ForgotPasswordRequest{Email: "invalid-email"}, handler.ForgotPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "invalid email format")
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/forgot-password", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ForgotPassword(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	t.Run("should return no content when password is reset successfully", func(t *testing.T) {
		svc, mockRepo, _, mockResetTokenSvc, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockResetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).Return("user-123", nil).Once()

		email, _ := domain.NewEmail("user@example.com")
		password, _ := domain.NewPassword("OldPass123!")
		user := domain.NewUser(email, password)
		user.ID = "user-123"

		mockRepo.EXPECT().FindByID(mock.Anything, "user-123").Return(user, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()
		mockResetTokenSvc.EXPECT().Invalidate(mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := nethttp.Request(
			t.Context(), http.MethodPost, "/auth/reset-password",
			map[string]any{"token": "valid-token", "password": "NewPass123!"}, handler.ResetPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should return bad request when token is invalid", func(t *testing.T) {
		svc, _, _, mockResetTokenSvc, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockResetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).Return("", domain.ErrInvalidResetToken).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/reset-password",
			map[string]any{"token": "invalid-token", "password": "NewPass123!"}, handler.ResetPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid reset token", got.Message)
	})

	t.Run("should return gone when token is expired", func(t *testing.T) {
		svc, _, _, mockResetTokenSvc, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockResetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).Return("", domain.ErrExpiredResetToken).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/reset-password",
			map[string]any{"token": "expired-token", "password": "NewPass123!"}, handler.ResetPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusGone, resp.StatusCode)
		assert.Equal(t, "reset token expired", got.Message)
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/reset-password", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ResetPassword(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})
}

func TestAuthHandler_ValidateResetToken(t *testing.T) {
	t.Run("should return ok when token is valid", func(t *testing.T) {
		svc, _, _, mockResetTokenSvc, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockResetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).Return("user-123", nil).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, ports.ValidateResetTokenResponse](
			t.Context(), http.MethodGet, "/auth/reset-password/validate?token=valid-token", nil, handler.ValidateResetToken)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, got.Valid)
	})

	t.Run("should return bad request when token is invalid", func(t *testing.T) {
		svc, _, _, mockResetTokenSvc, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockResetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).Return("", domain.ErrInvalidResetToken).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodGet, "/auth/reset-password/validate?token=invalid-token", nil, handler.ValidateResetToken)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid reset token", got.Message)
	})

	t.Run("should return gone when token is expired", func(t *testing.T) {
		svc, _, _, mockResetTokenSvc, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockResetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).Return("", domain.ErrExpiredResetToken).Once()

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodGet, "/auth/reset-password/validate?token=expired-token", nil, handler.ValidateResetToken)
		require.NoError(t, err)

		assert.Equal(t, http.StatusGone, resp.StatusCode)
		assert.Equal(t, "reset token expired", got.Message)
	})

	t.Run("should return bad request when token is empty", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodGet, "/auth/reset-password/validate?token=", nil, handler.ValidateResetToken)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "token is required", got.Message)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	t.Run("should return 201 when user is created successfully", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()
		tokenService.EXPECT().Store(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		emailSender.EXPECT().SendVerificationEmail(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.RegisterRequest, ports.RegisterResponse](
			t.Context(), http.MethodPost, "/auth/register",
			ports.RegisterRequest{Email: "user@example.com", Password: "ValidPass123!"}, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, "user@example.com", got.Email)
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/register", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.Register(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return 400 when email is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[ports.RegisterRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register",
			ports.RegisterRequest{Email: "invalid-email", Password: "ValidPass123!"}, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "invalid email format")
	})

	t.Run("should return 201 when user already exists for enumeration prevention", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(domain.ErrUserAlreadyExists).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.RegisterRequest, ports.RegisterResponse](
			t.Context(), http.MethodPost, "/auth/register",
			ports.RegisterRequest{Email: "user@example.com", Password: "ValidPass123!"}, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, got.ID)
	})

	t.Run("should return 500 when service fails with unknown error", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(assert.AnError).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.RegisterRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/register",
			ports.RegisterRequest{Email: "user@example.com", Password: "ValidPass123!"}, handler.Register)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("should return 200 when login is successful", func(t *testing.T) {
		svc, mockRepo, mockSecureTokenSvc, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)
		password, err := domain.NewPassword("ValidPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)

		mockRepo.EXPECT().FindVerifiedByEmail(mock.Anything, email).Return(user, nil).Once()
		mockSecureTokenSvc.EXPECT().Generate(mock.Anything, mock.Anything, mock.Anything).Return("refresh-token", nil).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.LoginRequest, ports.LoginResponse](
			t.Context(), http.MethodPost, "/auth/login",
			ports.LoginRequest{Email: "user@example.com", Password: "ValidPass123!"}, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, got.Token)
		assert.Equal(t, "refresh-token", got.RefreshToken)
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/login", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.Login(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return 400 when email is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[ports.LoginRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login",
			ports.LoginRequest{Email: "invalid-email", Password: "ValidPass123!"}, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "invalid email format")
	})

	t.Run("should return 401 when credentials are invalid", func(t *testing.T) {
		svc, mockRepo, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		mockRepo.EXPECT().FindVerifiedByEmail(mock.Anything, email).Return(domain.User{}, domain.ErrUserNotFound).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.LoginRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login",
			ports.LoginRequest{Email: "user@example.com", Password: "ValidPass123!"}, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "invalid credentials", got.Message)
	})

	t.Run("should return 500 when service fails with unknown error", func(t *testing.T) {
		svc, mockRepo, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)

		mockRepo.EXPECT().FindVerifiedByEmail(mock.Anything, email).Return(domain.User{}, assert.AnError).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.LoginRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/login",
			ports.LoginRequest{Email: "user@example.com", Password: "ValidPass123!"}, handler.Login)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	t.Run("should return 200 when refresh is successful", func(t *testing.T) {
		svc, _, mockSecureTokenSvc, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockSecureTokenSvc.EXPECT().Refresh(mock.Anything, "valid-refresh-token").
			Return("user-123", "user@example.com", "new-refresh-token", nil).
			Once()

		resp, got, err := nethttp.RequestWithResponse[ports.RefreshRequest, ports.RefreshResponse](
			t.Context(), http.MethodPost, "/auth/refresh",
			ports.RefreshRequest{RefreshToken: "valid-refresh-token"}, handler.Refresh)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, got.Token)
		assert.Equal(t, "new-refresh-token", got.RefreshToken)
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/refresh", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.Refresh(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return 400 when refresh token is empty", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[ports.RefreshRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/refresh",
			ports.RefreshRequest{RefreshToken: ""}, handler.Refresh)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "refresh_token is required", got.Message)
	})

	t.Run("should return 500 when service fails with unknown error", func(t *testing.T) {
		svc, _, mockSecureTokenSvc, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockSecureTokenSvc.EXPECT().Refresh(mock.Anything, "token").
			Return("", "", "", assert.AnError).
			Once()

		resp, got, err := nethttp.RequestWithResponse[ports.RefreshRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/refresh",
			ports.RefreshRequest{RefreshToken: "token"}, handler.Refresh)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	t.Run("should return 204 when logout is successful", func(t *testing.T) {
		svc, _, mockSecureTokenSvc, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockSecureTokenSvc.EXPECT().Invalidate(mock.Anything, "refresh-token").Return(nil).Once()

		resp, err := nethttp.Request(
			t.Context(), http.MethodPost, "/auth/logout",
			ports.LogoutRequest{RefreshToken: "refresh-token"}, handler.Logout)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.Logout(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return 400 when refresh token is empty", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[ports.LogoutRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/logout",
			ports.LogoutRequest{RefreshToken: ""}, handler.Logout)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "refresh_token is required", got.Message)
	})

	t.Run("should return 204 even when token is already invalid", func(t *testing.T) {
		svc, _, mockSecureTokenSvc, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		mockSecureTokenSvc.EXPECT().Invalidate(mock.Anything, "invalid-token").Return(domain.ErrRefreshTokenNotFound).Once()

		resp, err := nethttp.Request(
			t.Context(), http.MethodPost, "/auth/logout",
			ports.LogoutRequest{RefreshToken: "invalid-token"}, handler.Logout)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestAuthHandler_VerifyEmail(t *testing.T) {
	t.Run("should return 200 when email is verified successfully", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)
		password, err := domain.NewPassword("ValidPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)
		user.ID = "user-123"

		tokenService.EXPECT().Validate(mock.Anything, "valid-token").Return("user-123", nil).Once()
		userRepoMock.EXPECT().FindByID(mock.Anything, "user-123").Return(user, nil).Once()
		userRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()
		tokenService.EXPECT().Invalidate(mock.Anything, "valid-token").Return(nil).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.VerifyEmailRequest, ports.VerifyEmailResponse](
			t.Context(), http.MethodPost, "/auth/verify-email",
			ports.VerifyEmailRequest{Token: "valid-token"}, handler.VerifyEmail)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "Email verified successfully", got.Message)
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/verify-email", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.VerifyEmail(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})

	t.Run("should return 400 when token is empty", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[ports.VerifyEmailRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/verify-email",
			ports.VerifyEmailRequest{Token: ""}, handler.VerifyEmail)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "token is required", got.Message)
	})

	t.Run("should return 400 when verification token is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		tokenService.EXPECT().Validate(mock.Anything, "invalid-token").Return("", domain.ErrInvalidVerificationToken).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.VerifyEmailRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/verify-email",
			ports.VerifyEmailRequest{Token: "invalid-token"}, handler.VerifyEmail)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid verification token", got.Message)
	})

	t.Run("should return 410 when verification token is expired", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		tokenService.EXPECT().Validate(mock.Anything, "expired-token").Return("", domain.ErrExpiredVerificationToken).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.VerifyEmailRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/verify-email",
			ports.VerifyEmailRequest{Token: "expired-token"}, handler.VerifyEmail)
		require.NoError(t, err)

		assert.Equal(t, http.StatusGone, resp.StatusCode)
		assert.Equal(t, "verification token expired", got.Message)
	})

	t.Run("should return 500 when service fails with unknown error", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		tokenService := mocks.NewMockTokenService(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize, tokenService, verificationTokenTTL, verificationTokenSize, verificationURL)
		handler := ports.NewAuthHandler(authService)

		tokenService.EXPECT().Validate(mock.Anything, "token").Return("", assert.AnError).Once()

		resp, got, err := nethttp.RequestWithResponse[ports.VerifyEmailRequest, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/verify-email",
			ports.VerifyEmailRequest{Token: "token"}, handler.VerifyEmail)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestMapErrorToHTTPStatus(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		wantCode int
		wantMsg  string
	}{
		{"user already exists", domain.ErrUserAlreadyExists, http.StatusConflict, domain.ErrUserAlreadyExists.Error()},
		{"email empty", domain.ErrEmailIsEmpty, http.StatusBadRequest, domain.ErrEmailIsEmpty.Error()},
		{"invalid email", domain.ErrInvalidEmail, http.StatusBadRequest, domain.ErrInvalidEmail.Error()},
		{"password too short", domain.ErrPasswordTooShort, http.StatusBadRequest, domain.ErrPasswordTooShort.Error()},
		{"password too long", domain.ErrPasswordTooLong, http.StatusBadRequest, domain.ErrPasswordTooLong.Error()},
		{"password missing upper", domain.ErrPasswordMissingUpper, http.StatusBadRequest, domain.ErrPasswordMissingUpper.Error()},
		{"password missing lower", domain.ErrPasswordMissingLower, http.StatusBadRequest, domain.ErrPasswordMissingLower.Error()},
		{"password missing number", domain.ErrPasswordMissingNumber, http.StatusBadRequest, domain.ErrPasswordMissingNumber.Error()},
		{"password missing symbol", domain.ErrPasswordMissingSymbol, http.StatusBadRequest, domain.ErrPasswordMissingSymbol.Error()},
		{"invalid credentials", domain.ErrInvalidCredentials, http.StatusUnauthorized, domain.ErrInvalidCredentials.Error()},
		{"invalid refresh token", domain.ErrRefreshTokenInvalid, http.StatusUnauthorized, domain.ErrRefreshTokenInvalid.Error()},
		{"refresh token not found", domain.ErrRefreshTokenNotFound, http.StatusUnauthorized, domain.ErrRefreshTokenInvalid.Error()},
		{"refresh token expired", domain.ErrRefreshTokenExpired, http.StatusUnauthorized, domain.ErrRefreshTokenInvalid.Error()},
		{"invalid reset token", domain.ErrInvalidResetToken, http.StatusBadRequest, domain.ErrInvalidResetToken.Error()},
		{"expired reset token", domain.ErrExpiredResetToken, http.StatusGone, domain.ErrExpiredResetToken.Error()},
		{"invalid verification token", domain.ErrInvalidVerificationToken, http.StatusBadRequest, domain.ErrInvalidVerificationToken.Error()},
		{"expired verification token", domain.ErrExpiredVerificationToken, http.StatusGone, domain.ErrExpiredVerificationToken.Error()},
		{"unknown error", assert.AnError, http.StatusInternalServerError, "internal server error"},
	}

	for _, tc := range testCases {
		t.Run("should handle "+tc.name, func(t *testing.T) {
			status, message := ports.MapErrorToHTTPStatus(tc.err)
			assert.Equal(t, tc.wantCode, status)
			assert.Equal(t, tc.wantMsg, message)
		})
	}

	t.Run("should not leak internal wrapper text on unknown error", func(t *testing.T) {
		wrapped := fmt.Errorf("s.userRepository.Create: %w", assert.AnError)
		code, msg := ports.MapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, "internal server error", msg)
		assert.NotContains(t, msg, "userRepository")
		assert.NotContains(t, msg, "Create")
	})

	t.Run("should use canonical domain message when sentinel is wrapped", func(t *testing.T) {
		wrapped := fmt.Errorf("app.Register: %w", domain.ErrInvalidEmail)
		code, msg := ports.MapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, domain.ErrInvalidEmail.Error(), msg)
		assert.NotContains(t, msg, "app.Register")
	})
}
