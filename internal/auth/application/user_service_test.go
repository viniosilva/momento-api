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

func TestNewUserService(t *testing.T) {
	t.Run("should create user service", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		assert.NotNil(t, userService)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	defaultUserInput := application.UserInput{
		Email:    "user@example.com",
		Password: "ValidPass123!",
	}

	t.Run("should create user successfully", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		userRepoMock.EXPECT().HasByEmail(mock.Anything, mock.Anything).Return(false, nil).Once()
		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := userService.CreateUser(t.Context(), defaultUserInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, domain.Email("user@example.com"), got.Email)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		input := defaultUserInput
		input.Email = "invalid-email"

		_, err := userService.CreateUser(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		input := defaultUserInput
		input.Password = "invalid"

		_, err := userService.CreateUser(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		userRepoMock.EXPECT().HasByEmail(mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		_, err := userService.CreateUser(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("should return error when HasByEmail fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		userRepoMock.EXPECT().HasByEmail(mock.Anything, mock.Anything).
			Return(false, assert.AnError).
			Once()

		_, err := userService.CreateUser(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		userRepoMock := mocks.NewMockUserRepository(t)
		userService := application.NewUserService(userRepoMock)

		userRepoMock.EXPECT().HasByEmail(mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		userRepoMock.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		_, err := userService.CreateUser(t.Context(), defaultUserInput)

		assert.ErrorIs(t, err, assert.AnError)
	})
}
