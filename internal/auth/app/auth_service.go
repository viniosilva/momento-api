package app

import (
	"context"
	"errors"
	"fmt"

	"momento/internal/auth/domain"
)

type authService struct {
	userRepository     UserRepository
	jwtService         JWTService
	secureTokenService SecureTokenService
}

func NewAuthService(
	userRepository UserRepository,
	jwtService JWTService,
	secureTokenService SecureTokenService,
) *authService {
	return &authService{
		userRepository:     userRepository,
		jwtService:         jwtService,
		secureTokenService: secureTokenService,
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
		ID:        user.ID,
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

	token, err := s.jwtService.Generate(user.ID, user.Email)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("s.jwtService.Generate: %w", err)
	}

	refreshToken, err := s.secureTokenService.Generate(ctx, user.ID, string(user.Email))
	if err != nil {
		return LoginOutput{}, fmt.Errorf("s.secureTokenService.Generate: %w", err)
	}

	return LoginOutput{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, input RefreshTokenInput) (LoginOutput, error) {
	userID, email, newRefreshToken, err := s.secureTokenService.Refresh(ctx, input.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) || errors.Is(err, domain.ErrRefreshTokenExpired) {
			return LoginOutput{}, domain.ErrRefreshTokenInvalid
		}
		return LoginOutput{}, fmt.Errorf("s.secureTokenService.Refresh: %w", err)
	}

	token, err := s.jwtService.Generate(userID, domain.Email(email))
	if err != nil {
		return LoginOutput{}, fmt.Errorf("s.jwtService.Generate: %w", err)
	}

	return LoginOutput{
		Token:        token,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, input LogoutInput) error {
	if err := s.secureTokenService.Invalidate(ctx, input.RefreshToken); err != nil {
		// Even if token doesn't exist or is already invalid, we return nil
		// This is intentional - logout should succeed even if token is already invalid
		return nil
	}

	return nil
}
