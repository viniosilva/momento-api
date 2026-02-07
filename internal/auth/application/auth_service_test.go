package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pinnado/internal/auth/application"
	"pinnado/internal/auth/domain"
	"pinnado/mocks"
)

func TestNewAuthService(t *testing.T) {
	t.Run("should create auth service", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(userRepoMock)

		assert.NotNil(t, authService)
	})
}

func TestAuthService_Register(t *testing.T) {
	defaultUserInput := application.UserInput{
		Email:    "user@example.com",
		Password: "ValidPass123!",
	}

	t.Run("should create user successfully", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewAuthService(userRepoMock)

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
		userService := application.NewAuthService(userRepoMock)

		input := defaultUserInput
		input.Email = "invalid-email"

		_, err := userService.Register(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewAuthService(userRepoMock)

		input := defaultUserInput
		input.Password = "short"

		_, err := userService.Register(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewAuthService(userRepoMock)

		userRepoMock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		_, err := userService.Register(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("should return error when ExistsByEmail fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewAuthService(userRepoMock)

		userRepoMock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).
			Return(false, assert.AnError).
			Once()

		_, err := userService.Register(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewAuthService(userRepoMock)

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

func TestAuthService_Login(t *testing.T) {
	defaultLoginInput := application.LoginInput{
		Email:    "user@example.com",
		Password: "ValidPass123!",
	}

	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(userRepoMock)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		password, err := domain.NewPassword(defaultLoginInput.Password)
		require.NoError(t, err)

		user := domain.NewUser(email, password)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(user, nil).
			Once()

		got, err := authService.Login(t.Context(), defaultLoginInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, email, got.Email)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
	})

	t.Run("should return error when credentials are invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		authService := application.NewAuthService(userRepoMock)

		email, err := domain.NewEmail(defaultLoginInput.Email)
		require.NoError(t, err)

		userRepoMock.EXPECT().FindByEmail(mock.Anything, email).
			Return(domain.User{}, domain.ErrUserNotFound).
			Once()

		_, err = authService.Login(t.Context(), defaultLoginInput)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})
}
