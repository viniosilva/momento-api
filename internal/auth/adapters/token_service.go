package adapters

import (
	"context"
	"fmt"
	"time"

	"momento/internal/auth/domain"

	"github.com/redis/go-redis/v9"
)

const VerificationTokenPrefix = "verification_token:"

type tokenService struct {
	redis  *redis.Client
	prefix string
}

func NewTokenService(redis *redis.Client, prefix string) *tokenService {
	return &tokenService{
		redis:  redis,
		prefix: prefix,
	}
}

func (s *tokenService) Store(ctx context.Context, token string, userID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", s.prefix, token)
	if err := s.redis.Set(ctx, key, userID, ttl).Err(); err != nil {
		return fmt.Errorf("s.redis.Set: %w", err)
	}

	return nil
}

func (s *tokenService) Validate(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s%s", s.prefix, token)
	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrInvalidVerificationToken
	}
	if err != nil {
		return "", fmt.Errorf("s.redis.Get: %w", err)
	}

	return val, nil
}

func (s *tokenService) Invalidate(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s%s", s.prefix, token)
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("s.redis.Del: %w", err)
	}

	return nil
}
