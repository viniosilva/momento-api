package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"momento/internal/auth/domain"
)

type authService struct {
	userRepository     UserRepository
	jwtService         JWTService
	secureTokenService SecureTokenService
	resetTokenService  ResetTokenService
	emailSender        EmailSender
	resetTokenTTL      time.Duration
	resetTokenSize     int
}

func NewAuthService(
	userRepository UserRepository,
	jwtService JWTService,
	secureTokenService SecureTokenService,
	resetTokenService ResetTokenService,
	emailSender EmailSender,
	resetTokenTTL time.Duration,
	resetTokenSize int,
) *authService {
	return &authService{
		userRepository:     userRepository,
		jwtService:         jwtService,
		secureTokenService: secureTokenService,
		resetTokenService:  resetTokenService,
		emailSender:        emailSender,
		resetTokenTTL:      resetTokenTTL,
		resetTokenSize:     resetTokenSize,
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

func (s *authService) ForgotPassword(ctx context.Context, input ForgotPasswordInput) error {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return err
	}

	token, err := domain.NewResetToken(s.resetTokenSize)
	if err != nil {
		return fmt.Errorf("domain.NewResetToken: %w", err)
	}

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil // We don't want to reveal whether the email exists or not, so we return nil even if user is not found
	}

	if err := s.resetTokenService.Store(ctx, token, user.ID, s.resetTokenTTL); err != nil {
		return fmt.Errorf("s.resetTokenSvc.Store: %w", err)
	}

	if err := s.emailSender.SendResetPasswordEmail(ctx, email.String(), token.String()); err != nil {
		return fmt.Errorf("s.emailSender.SendResetPasswordEmail: %w", err)
	}

	return nil
}

func (s *authService) ValidateResetToken(ctx context.Context, input ValidateResetTokenInput) (string, error) {
	token := domain.ResetToken(input.Token)
	return s.resetTokenService.Validate(ctx, token)
}

func (s *authService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	token := domain.ResetToken(input.Token)
	userID, err := s.resetTokenService.Validate(ctx, token)
	if err != nil {
		return fmt.Errorf("s.resetTokenService.Validate: %w", err)
	}

	password, err := domain.NewPassword(input.Password)
	if err != nil {
		return fmt.Errorf("domain.NewPassword: %w", err)
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("s.userRepository.FindByID: %w", err)
	}

	user.UpdatePassword(password)

	if err := s.userRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("s.userRepository.Update: %w", err)
	}

	if err := s.resetTokenService.Invalidate(ctx, token); err != nil {
		return fmt.Errorf("s.resetTokenService.Invalidate: %w", err)
	}

	return nil
}
