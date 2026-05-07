package adapters

import (
	"context"
	"fmt"
	"time"

	"momento/internal/auth/domain"

	"github.com/redis/go-redis/v9"
)

const ResetTokenPrefix = "reset_token:"

type resetTokenService struct {
	redis *redis.Client
}

func NewResetTokenService(redis *redis.Client) *resetTokenService {
	return &resetTokenService{redis: redis}
}

func (s *resetTokenService) Store(ctx context.Context, token domain.ResetToken, userID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", ResetTokenPrefix, token.String())
	if err := s.redis.Set(ctx, key, userID, ttl).Err(); err != nil {
		return fmt.Errorf("s.redis.Set: %w", err)
	}

	return nil
}

func (s *resetTokenService) Validate(ctx context.Context, token domain.ResetToken) (string, error) {
	key := fmt.Sprintf("%s%s", ResetTokenPrefix, token.String())
	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrInvalidResetToken
	}
	if err != nil {
		return "", fmt.Errorf("s.redis.Get: %w", err)
	}

	return val, nil
}

func (s *resetTokenService) Invalidate(ctx context.Context, token domain.ResetToken) error {
	key := fmt.Sprintf("%s%s", ResetTokenPrefix, token.String())
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("s.redis.Del: %w", err)
	}

	return nil
}
