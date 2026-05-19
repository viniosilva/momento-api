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
		assert.Equal(t, "localhost", got.PG.Host)
		assert.Equal(t, "5432", got.PG.Port)
		assert.Equal(t, "momento", got.PG.User)
		assert.Equal(t, "momento", got.PG.Pass)
		assert.Equal(t, "momento", got.PG.DBName)
		assert.Equal(t, "disable", got.PG.SSLMode)
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
		assert.Equal(t, "localhost", got.PG.Host)
		assert.Equal(t, "5432", got.PG.Port)
		assert.Equal(t, "momento", got.PG.User)
		assert.Equal(t, "momento", got.PG.Pass)
		assert.Equal(t, "momento", got.PG.DBName)
		assert.Equal(t, "disable", got.PG.SSLMode)
		assert.Equal(t, 10*time.Second, got.PG.ConnectTimeout)
		assert.Equal(t, "your-secret-key-change-in-production", got.JWT.Secret)
		assert.Equal(t, 12*time.Hour, got.JWT.Expiration)
		assert.Equal(t, 7*24*time.Hour, got.JWT.RefreshTokenExpiration)
		assert.Equal(t, "localhost", got.Redis.Host)
		assert.Equal(t, "6379", got.Redis.Port)
		assert.Equal(t, "", got.Redis.Pass)
		assert.Equal(t, 0, got.Redis.DB)
	})
}

var envVars []string = []string{
	"API_HOST",
	"API_PORT",
	"PG_HOST",
	"PG_PORT",
	"PG_USER",
	"PG_PASS",
	"PG_DBNAME",
	"PG_SSLMODE",
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
