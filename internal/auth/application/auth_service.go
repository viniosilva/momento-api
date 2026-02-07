package application

import (
	"context"
	"errors"
	"fmt"
	"pinnado/internal/auth/domain"
)

type authService struct {
	userRepository UserRepository
}

func NewAuthService(userRepository UserRepository) *authService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) Register(ctx context.Context, input UserInput) (UserOutput, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return UserOutput{}, err
	}

	password, err := domain.NewPassword(input.Password)
	if err != nil {
		return UserOutput{}, err
	}

	exists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return UserOutput{}, fmt.Errorf("s.userRepository.ExistsByEmail: %w", err)
	}
	if exists {
		return UserOutput{}, domain.ErrUserAlreadyExists
	}

	user := domain.NewUser(email, password)
	if err := s.userRepository.Create(ctx, user); err != nil {
		return UserOutput{}, fmt.Errorf("s.userRepository.Create: %w", err)
	}

	return UserOutput{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return LoginOutput{}, err
	}

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return LoginOutput{}, domain.ErrInvalidCredentials
		}
		return LoginOutput{}, fmt.Errorf("s.userRepository.FindByEmail: %w", err)
	}

	if err := user.Password.Compare(input.Password); err != nil {
		return LoginOutput{}, domain.ErrInvalidCredentials
	}

	return LoginOutput{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
