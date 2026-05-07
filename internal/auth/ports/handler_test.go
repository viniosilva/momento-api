package ports_test

import (
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
	resetTokenSize = 32
	resetTokenTTL  = 1 * time.Hour
)

func newAuthServiceWithMocks(t *testing.T) (ports.AuthService, *mocks.MockUserRepository, *mocks.MockSecureTokenService, *mocks.MockResetTokenService, *mocks.MockEmailSender) {
	t.Helper()
	userRepo := mocks.NewMockUserRepository(t)
	secureTokenSvc := mocks.NewMockSecureTokenService(t)
	jwtService := adapters.NewJWTService(secretTest, expirationTest)
	resetTokenSvc := mocks.NewMockResetTokenService(t)
	emailSender := mocks.NewMockEmailSender(t)

	svc := app.NewAuthService(userRepo, jwtService, secureTokenSvc, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)
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

		resp, err := nethttp.Request[map[string]any](
			t.Context(), http.MethodPost, "/auth/forgot-password",
			map[string]any{"email": "user@example.com"}, handler.ForgotPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should return no content even when user does not exist", func(t *testing.T) {
		svc, mockRepo, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		email, _ := domain.NewEmail("user@example.com")
		mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(domain.User{}, domain.ErrUserNotFound).Once()

		resp, err := nethttp.Request[map[string]any](
			t.Context(), http.MethodPost, "/auth/forgot-password",
			map[string]any{"email": "user@example.com"}, handler.ForgotPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		svc, _, _, _, _ := newAuthServiceWithMocks(t)
		handler := ports.NewAuthHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[map[string]any, nethttp.ErrorResponse](
			t.Context(), http.MethodPost, "/auth/forgot-password",
			map[string]any{"email": "invalid-email"}, handler.ForgotPassword)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "invalid email format")
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

		resp, err := nethttp.Request[map[string]any](
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
