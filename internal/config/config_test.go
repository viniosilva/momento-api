package config_test

import (
	"os"
	"path/filepath"
	"pinnado/internal/config"
	"testing"
	"time"

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
		assert.Equal(t, "localhost", got.Mongo.Host)
		assert.Equal(t, "27017", got.Mongo.Port)
		assert.Equal(t, "pinnado", got.Mongo.DBName)
		assert.Equal(t, "admin", got.Mongo.User)
		assert.Equal(t, "admin", got.Mongo.Pass)
		assert.Equal(t, 3, got.Mongo.MaxRetries)
		assert.Equal(t, 2*time.Second, got.Mongo.RetryDelay)
		assert.Equal(t, 10*time.Second, got.Mongo.ConnectTimeout)
		assert.Equal(t, "your-secret-key-change-in-production", got.JWT.Secret)
		assert.Equal(t, 12*time.Hour, got.JWT.Expiration)
	})

	t.Run("should load values from .env.example file", func(t *testing.T) {
		t.Cleanup(func() {
			clearEnvVars(t)
		})

		envExamplePath := filepath.Join("..", "..", ".env.example")
		got := config.LoadConfig(config.WithCustomPath(envExamplePath))

		assert.Equal(t, "0.0.0.0", got.Api.Host)
		assert.Equal(t, "8080", got.Api.Port)
		assert.Equal(t, "localhost", got.Mongo.Host)
		assert.Equal(t, "27017", got.Mongo.Port)
		assert.Equal(t, "pinnado", got.Mongo.DBName)
		assert.Equal(t, "admin", got.Mongo.User)
		assert.Equal(t, "admin", got.Mongo.Pass)
		assert.Equal(t, 3, got.Mongo.MaxRetries)
		assert.Equal(t, 2*time.Second, got.Mongo.RetryDelay)
		assert.Equal(t, 10*time.Second, got.Mongo.ConnectTimeout)
		assert.Equal(t, "your-secret-key-change-in-production", got.JWT.Secret)
		assert.Equal(t, 12*time.Hour, got.JWT.Expiration)
	})

	t.Run("should return default max retries when conversion fails", func(t *testing.T) {
		t.Cleanup(func() {
			clearEnvVars(t)
		})

		os.Setenv("MONGO_MAX_RETRIES", "invalid")

		got := config.LoadConfig()

		assert.Equal(t, 3, got.Mongo.MaxRetries)
	})
}

var envVars []string = []string{
	"API_HOST",
	"API_PORT",
	"MONGO_HOST",
	"MONGO_PORT",
	"MONGO_DB",
	"MONGO_USER",
	"MONGO_PASS",
	"MONGO_MAX_RETRIES",
	"MONGO_RETRY_DELAY_MS",
	"MONGO_CONNECT_TIMEOUT_MS",
	"JWT_SECRET",
	"JWT_EXPIRATION_MS",
}

func clearEnvVars(t *testing.T) {
	t.Helper()

	for _, key := range envVars {
		os.Unsetenv(key)
	}
}
