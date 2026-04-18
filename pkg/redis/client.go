package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, host, port, pass string, db int) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return cli, nil
}
