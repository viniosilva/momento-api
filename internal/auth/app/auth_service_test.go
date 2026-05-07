package app_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"momento/internal/auth/adapters"
	"momento/internal/auth/app"
	"momento/internal/auth/domain"
	"momento/internal/auth/mocks"
)

const (
	secretTest     = "secretTest"
	expirationTest = 5 * time.Second
	resetTokenSize = 32
	resetTokenTTL  = 1 * time.Hour
)

func TestNewAuthService(t *testing.T) {
	t.Run("should create auth service", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)

		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		assert.NotNil(t, authService)
	})
}

func TestAuthService_Register(t *testing.T) {
	defaultUserInput := app.UserInput{
		Email:    "user@example.com",
		Password: "ValidPass123!",
	}

	t.Run("should create user successfully", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		userService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		userRepoMock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, nil).Once()
		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := userService.Register(t.Context(), defaultUserInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, domain.Email("user@example.com"), got.Email)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		userService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		input := defaultUserInput
		input.Email = "invalid-email"

		_, err := userService.Register(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		userService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		input := defaultUserInput
		input.Password = "short"

		_, err := userService.Register(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		userService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		userRepoMock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		_, err := userService.Register(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("should return error when ExistsByEmail fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		userService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		userRepoMock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).
			Return(false, assert.AnError).
			Once()

		_, err := userService.Register(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		userService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		userRepoMock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		_, err := userService.Register(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestAuthService_Logout(t *testing.T) {
	const existingToken = "existing-refresh-token"

	t.Run("should logout successfully by invalidating refresh token", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Invalidate(mock.Anything, existingToken).
			Return(nil).
			Once()

		err := authService.Logout(t.Context(), app.LogoutInput{RefreshToken: existingToken})
		require.NoError(t, err)
	})

	t.Run("should return nil even when token is already invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Invalidate(mock.Anything, existingToken).
			Return(domain.ErrRefreshTokenNotFound).
			Once()

		err := authService.Logout(t.Context(), app.LogoutInput{RefreshToken: existingToken})
		require.NoError(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	defaultLoginInput := app.LoginInput{
		Email:    "user@example.com",
		Password: "ValidPass123!",
	}

	t.Run("should login successfully and return token and refresh token", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword(defaultLoginInput.Password)
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()
		tokenSvcMock.EXPECT().Generate(mock.Anything, user.ID, string(user.Email)).
			Return("refresh-token-abc", nil).
			Once()

		got, err := authService.Login(t.Context(), defaultLoginInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.Token)
		assert.Equal(t, "refresh-token-abc", got.RefreshToken)
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		input := defaultLoginInput
		input.Email = "invalid-email"

		_, err := authService.Login(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return error when credentials are invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, domain.ErrUserNotFound).
			Once()

		_, err = authService.Login(t.Context(), defaultLoginInput)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("should return error when password is incorrect", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword("OtherPass123!")
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()

		_, err = authService.Login(t.Context(), defaultLoginInput)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("should return error when FindByEmail fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, assert.AnError).
			Once()

		_, err = authService.Login(t.Context(), defaultLoginInput)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when jwtService.Generate fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtMock := mocks.NewMockJWTService(t)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtMock, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword(defaultLoginInput.Password)
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()
		jwtMock.EXPECT().Generate(mock.Anything, mock.Anything).
			Return("", assert.AnError).
			Once()

		_, err = authService.Login(t.Context(), defaultLoginInput)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when secureTokenService.Generate fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword(defaultLoginInput.Password)
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()
		tokenSvcMock.EXPECT().Generate(mock.Anything, mock.Anything, mock.Anything).
			Return("", assert.AnError).
			Once()

		_, err = authService.Login(t.Context(), defaultLoginInput)

		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	const existingToken = "existing-refresh-token"

	t.Run("should rotate refresh token and return new tokens", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Refresh(mock.Anything, existingToken).
			Return("user-123", "user@example.com", "new-refresh-token", nil).
			Once()

		got, err := authService.RefreshToken(t.Context(), app.RefreshTokenInput{RefreshToken: existingToken})
		require.NoError(t, err)

		assert.NotEmpty(t, got.Token)
		assert.Equal(t, "new-refresh-token", got.RefreshToken)
	})

	t.Run("should return ErrRefreshTokenInvalid when token is not found", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Refresh(mock.Anything, existingToken).
			Return("", "", "", domain.ErrRefreshTokenNotFound).
			Once()

		_, err := authService.RefreshToken(t.Context(), app.RefreshTokenInput{RefreshToken: existingToken})

		assert.ErrorIs(t, err, domain.ErrRefreshTokenInvalid)
	})

	t.Run("should return ErrRefreshTokenInvalid when token is expired", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Refresh(mock.Anything, existingToken).
			Return("", "", "", domain.ErrRefreshTokenExpired).
			Once()

		_, err := authService.RefreshToken(t.Context(), app.RefreshTokenInput{RefreshToken: existingToken})

		assert.ErrorIs(t, err, domain.ErrRefreshTokenInvalid)
	})

	t.Run("should return error when secureTokenService.Refresh fails with unknown error", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Refresh(mock.Anything, existingToken).
			Return("", "", "", assert.AnError).
			Once()

		_, err := authService.RefreshToken(t.Context(), app.RefreshTokenInput{RefreshToken: existingToken})

		assert.ErrorIs(t, err, assert.AnError)
		assert.NotErrorIs(t, err, domain.ErrRefreshTokenInvalid)
	})

	t.Run("should return error when jwtService.Generate fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtMock := mocks.NewMockJWTService(t)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtMock, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		tokenSvcMock.EXPECT().Refresh(mock.Anything, existingToken).
			Return("user-123", "user@example.com", "new-refresh-token", nil).
			Once()
		jwtMock.EXPECT().Generate(mock.Anything, mock.Anything).
			Return("", assert.AnError).
			Once()

		_, err := authService.RefreshToken(t.Context(), app.RefreshTokenInput{RefreshToken: existingToken})

		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestAuthService_ForgotPassword(t *testing.T) {
	defaultInput := app.ForgotPasswordInput{
		Email: "user@example.com",
	}

	t.Run("should send reset email when user exists", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword("ValidPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()

		resetTokenSvc.EXPECT().Store(mock.Anything, mock.Anything, user.ID, resetTokenTTL).
			Return(nil).
			Once()

		emailSender.EXPECT().SendResetPasswordEmail(mock.Anything, defaultInput.Email, mock.Anything).
			Return(nil).
			Once()

		err = authService.ForgotPassword(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return nil when user does not exist (security)", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultInput.Email)
		require.NoError(t, err)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, domain.ErrUserNotFound).
			Once()

		err = authService.ForgotPassword(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		input := app.ForgotPasswordInput{Email: "invalid-email"}
		err := authService.ForgotPassword(t.Context(), input)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return nil even when FindByEmail fails (security)", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultInput.Email)
		require.NoError(t, err)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, assert.AnError).
			Once()

		err = authService.ForgotPassword(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when reset token service Store fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword("ValidPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()

		resetTokenSvc.EXPECT().Store(mock.Anything, mock.Anything, user.ID, resetTokenTTL).
			Return(assert.AnError).
			Once()

		err = authService.ForgotPassword(t.Context(), defaultInput)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when email sender fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		email, err := domain.NewEmail(defaultInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword("ValidPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()

		resetTokenSvc.EXPECT().Store(mock.Anything, mock.Anything, user.ID, resetTokenTTL).
			Return(nil).
			Once()

		emailSender.EXPECT().SendResetPasswordEmail(mock.Anything, defaultInput.Email, mock.Anything).
			Return(assert.AnError).
			Once()

		err = authService.ForgotPassword(t.Context(), defaultInput)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestAuthService_ValidateResetToken(t *testing.T) {
	t.Run("should return user ID when token is valid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		userID, err := authService.ValidateResetToken(t.Context(), app.ValidateResetTokenInput{Token: "valid-token"})
		require.NoError(t, err)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("", domain.ErrInvalidResetToken).
			Once()

		_, err := authService.ValidateResetToken(t.Context(), app.ValidateResetTokenInput{Token: "invalid-token"})
		assert.ErrorIs(t, err, domain.ErrInvalidResetToken)
	})

	t.Run("should return error when token is expired", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("", domain.ErrExpiredResetToken).
			Once()

		_, err := authService.ValidateResetToken(t.Context(), app.ValidateResetTokenInput{Token: "expired-token"})
		assert.ErrorIs(t, err, domain.ErrExpiredResetToken)
	})

	t.Run("should return error when validate fails with unknown error", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("", assert.AnError).
			Once()

		_, err := authService.ValidateResetToken(t.Context(), app.ValidateResetTokenInput{Token: "token"})
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestAuthService_ResetPassword(t *testing.T) {
	defaultInput := app.ResetPasswordInput{
		Token:    "valid-token",
		Password: "NewPass123!",
	}

	t.Run("should reset password successfully", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)
		password, err := domain.NewPassword("OldPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)
		user.ID = "user-123"

		userRepoMock.EXPECT().FindByID(mock.Anything, "user-123").
			Return(user, nil).
			Once()
		userRepoMock.EXPECT().Update(mock.Anything, mock.Anything).
			Return(nil).
			Once()
		resetTokenSvc.EXPECT().Invalidate(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		err = authService.ResetPassword(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("", domain.ErrInvalidResetToken).
			Once()

		err := authService.ResetPassword(t.Context(), defaultInput)
		assert.ErrorIs(t, err, domain.ErrInvalidResetToken)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		input := app.ResetPasswordInput{
			Token:    "valid-token",
			Password: "short",
		}
		err := authService.ResetPassword(t.Context(), input)
		assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
	})

	t.Run("should return error when FindByID fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		userRepoMock.EXPECT().FindByID(mock.Anything, "user-123").
			Return(domain.User{}, assert.AnError).
			Once()

		err := authService.ResetPassword(t.Context(), defaultInput)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when Update fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)
		password, err := domain.NewPassword("OldPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)
		user.ID = "user-123"

		userRepoMock.EXPECT().FindByID(mock.Anything, "user-123").
			Return(user, nil).
			Once()
		userRepoMock.EXPECT().Update(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		err = authService.ResetPassword(t.Context(), defaultInput)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when Invalidate fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		tokenSvcMock := mocks.NewMockSecureTokenService(t)
		jwtService := adapters.NewJWTService(secretTest, expirationTest)
		resetTokenSvc := mocks.NewMockResetTokenService(t)
		emailSender := mocks.NewMockEmailSender(t)
		authService := app.NewAuthService(userRepoMock, jwtService, tokenSvcMock, resetTokenSvc, emailSender, resetTokenTTL, resetTokenSize)

		resetTokenSvc.EXPECT().Validate(mock.Anything, mock.Anything).
			Return("user-123", nil).
			Once()

		email, err := domain.NewEmail("user@example.com")
		require.NoError(t, err)
		password, err := domain.NewPassword("OldPass123!")
		require.NoError(t, err)
		user := domain.NewUser(email, password)
		user.ID = "user-123"

		userRepoMock.EXPECT().FindByID(mock.Anything, "user-123").
			Return(user, nil).
			Once()
		userRepoMock.EXPECT().Update(mock.Anything, mock.Anything).
			Return(nil).
			Once()
		resetTokenSvc.EXPECT().Invalidate(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		err = authService.ResetPassword(t.Context(), defaultInput)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
