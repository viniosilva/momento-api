package application

import (
	"context"
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
