package application

import (
	"context"
	"errors"
	"fmt"
	"momento/internal/auth/domain"
)

type AuthService struct {
	userRepository UserRepository
	jwtService     JWTService
}

func NewAuthService(userRepository UserRepository, jwtService JWTService) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, input UserInput) (UserOutput, error) {
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

func (s *AuthService) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
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

	token, err := s.jwtService.Generate(user.ID.Hex(), user.Email)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("s.jwtService.Generate: %w", err)
	}

	return LoginOutput{
		Token: token,
	}, nil
}
