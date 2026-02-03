package application

import (
	"context"
	"fmt"
	"pinnado/internal/auth/domain"
)

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(ctx context.Context, input UserInput) (UserOutput, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return UserOutput{}, fmt.Errorf("domain.NewEmail: %w", err)
	}

	password, err := domain.NewPassword(input.Password)
	if err != nil {
		return UserOutput{}, fmt.Errorf("domain.NewPassword: %w", err)
	}

	exists, err := s.userRepository.HasByEmail(ctx, email)
	if err != nil {
		return UserOutput{}, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return UserOutput{}, domain.ErrUserAlreadyExists
	}

	user := domain.NewUser(email, password)
	if err := s.userRepository.Create(ctx, user); err != nil {
		return UserOutput{}, fmt.Errorf("failed to create user: %w", err)
	}

	return UserOutput{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
