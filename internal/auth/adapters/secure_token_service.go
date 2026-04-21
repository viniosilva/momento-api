package adapters

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"momento/internal/auth/domain"

	"github.com/redis/go-redis/v9"
)

const secureTokenBytes = 32

type secureTokenService struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func NewSecureTokenService(client redis.UniversalClient, ttl time.Duration) *secureTokenService {
	return &secureTokenService{
		client: client,
		ttl:    ttl,
	}
}

func (s *secureTokenService) Generate(ctx context.Context, userID, email string) (string, error) {
	token, err := generateSecureToken()
	if err != nil {
		return "", fmt.Errorf("generateSecureToken: %w", err)
	}

	if err := s.save(ctx, token, userID, email); err != nil {
		return "", err
	}

	return token, nil
}

func (s *secureTokenService) Refresh(ctx context.Context, token string) (userID, email, newToken string, err error) {
	userID, email, err = s.get(ctx, token)
	if err != nil {
		return
	}

	if err = s.delete(ctx, token); err != nil {
		return
	}

	newToken, err = s.Generate(ctx, userID, email)
	return
}

func (s *secureTokenService) Invalidate(ctx context.Context, token string) error {
	return s.delete(ctx, token)
}

func (s *secureTokenService) save(ctx context.Context, token, userID, email string) error {
	data := refreshTokenData{UserID: userID, Email: email}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err := s.client.Set(ctx, token, payload, s.ttl).Err(); err != nil {
		return fmt.Errorf("s.client.Set: %w", err)
	}

	return nil
}

func (s *secureTokenService) get(ctx context.Context, token string) (userID, email string, err error) {
	val, err := s.client.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", domain.ErrRefreshTokenNotFound
		}
		return "", "", fmt.Errorf("s.client.Get: %w", err)
	}

	var data refreshTokenData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return "", "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return data.UserID, data.Email, nil
}

func (s *secureTokenService) delete(ctx context.Context, token string) error {
	if err := s.client.Del(ctx, token).Err(); err != nil {
		return fmt.Errorf("s.client.Del: %w", err)
	}

	return nil
}

func generateSecureToken() (string, error) {
	b := make([]byte, secureTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
