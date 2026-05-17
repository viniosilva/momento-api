package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"momento/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("should return default values when env vars are not set", func(t *testing.T) {
		t.Cleanup(func() {
			clearEnvVars(t)
		})

		got := config.LoadConfig()

		assert.Equal(t, "", got.Api.Host)
		assert.Equal(t, "8080", got.Api.Port)
		assert.Equal(t, "postgres://momento:momento@localhost:5432/momento?sslmode=disable", got.PG.DSN)
		assert.Equal(t, 3, got.PG.MaxRetries)
		assert.Equal(t, 2*time.Second, got.PG.RetryDelay)
		assert.Equal(t, 10*time.Second, got.PG.ConnectTimeout)
		assert.Equal(t, "your-secret-key-change-in-production", got.JWT.Secret)
		assert.Equal(t, 12*time.Hour, got.JWT.Expiration)
		assert.Equal(t, 7*24*time.Hour, got.JWT.RefreshTokenExpiration)
		assert.Equal(t, "localhost", got.Redis.Host)
		assert.Equal(t, "6379", got.Redis.Port)
		assert.Equal(t, "", got.Redis.Pass)
		assert.Equal(t, 0, got.Redis.DB)
	})

	t.Run("should load values from .env.example file", func(t *testing.T) {
		t.Cleanup(func() {
			clearEnvVars(t)
		})

		envExamplePath := filepath.Join("..", "..", ".env.example")
		got := config.LoadConfig(config.WithCustomPath(envExamplePath))

		assert.Equal(t, "0.0.0.0", got.Api.Host)
		assert.Equal(t, "8080", got.Api.Port)
		assert.Equal(t, "postgres://momento:momento@localhost:5432/momento?sslmode=disable", got.PG.DSN)
		assert.Equal(t, 3, got.PG.MaxRetries)
		assert.Equal(t, 2*time.Second, got.PG.RetryDelay)
		assert.Equal(t, 10*time.Second, got.PG.ConnectTimeout)
		assert.Equal(t, "your-secret-key-change-in-production", got.JWT.Secret)
		assert.Equal(t, 12*time.Hour, got.JWT.Expiration)
		assert.Equal(t, 7*24*time.Hour, got.JWT.RefreshTokenExpiration)
		assert.Equal(t, "localhost", got.Redis.Host)
		assert.Equal(t, "6379", got.Redis.Port)
		assert.Equal(t, "", got.Redis.Pass)
		assert.Equal(t, 0, got.Redis.DB)
	})

	t.Run("should return default max retries when conversion fails", func(t *testing.T) {
		t.Cleanup(func() {
			clearEnvVars(t)
		})

		os.Setenv("PG_MAX_RETRIES", "invalid")

		got := config.LoadConfig()

		assert.Equal(t, 3, got.PG.MaxRetries)
	})
}

var envVars []string = []string{
	"API_HOST",
	"API_PORT",
	"DATABASE_URL",
	"PG_MAX_RETRIES",
	"PG_RETRY_DELAY_MS",
	"PG_CONNECT_TIMEOUT_MS",
	"JWT_SECRET",
	"JWT_EXPIRATION_MS",
	"REFRESH_TOKEN_EXPIRATION_MS",
	"REDIS_HOST",
	"REDIS_PORT",
	"REDIS_PASS",
	"REDIS_DB",
}

func clearEnvVars(t *testing.T) {
	t.Helper()

	for _, key := range envVars {
		os.Unsetenv(key)
	}
}
