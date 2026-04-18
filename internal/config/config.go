package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Api   ApiConfig
	Mongo MongoConfig
	JWT   JWTConfig
	Redis RedisConfig
}

type ApiConfig struct {
	Host string
	Port string
}

type MongoConfig struct {
	Host           string
	Port           string
	DBName         string
	User           string
	Pass           string
	MaxRetries     int
	RetryDelay     time.Duration
	ConnectTimeout time.Duration
}

type JWTConfig struct {
	Secret                 string
	Expiration             time.Duration
	RefreshTokenExpiration time.Duration
}

type RedisConfig struct {
	Host string
	Port string
	Pass string
	DB   int
}

const (
	defaultEnvPath                = ".env"
	defaultApiHost                = ""
	defaultApiPort                = "8080"
	defaultMongoHost              = "localhost"
	defaultMongoPort              = "27017"
	defaultMongoDB                = "momento"
	defaultMongoUser              = "admin"
	defaultMongoPass              = "admin"
	defaultMongoMaxRetries        = 3
	defaultMongoRetryDelay        = 2 * time.Second
	defaultMongoConnectTimeout    = 10 * time.Second
	defaultJWTSecret              = "your-secret-key-change-in-production"
	defaultJWTExpiration          = 12 * time.Hour
	defaultRefreshTokenExpiration = 7 * 24 * time.Hour
	defaultRedisHost              = "localhost"
	defaultRedisPort              = "6379"
	defaultRedisPass              = ""
	defaultRedisDB                = 0
)

type LoadConfigOption func(*loadConfigOptions)

type loadConfigOptions struct {
	envPath string
}

func WithCustomPath(path string) LoadConfigOption {
	return func(opts *loadConfigOptions) {
		opts.envPath = path
	}
}

func LoadConfig(opts ...LoadConfigOption) Config {
	options := &loadConfigOptions{
		envPath: defaultEnvPath,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := godotenv.Load(options.envPath); err != nil {
		log.Println("No .env file found, using default values")
	}

	return Config{
		Api: ApiConfig{
			Host: getEnv("API_HOST", defaultApiHost),
			Port: getEnv("API_PORT", defaultApiPort),
		},
		Mongo: MongoConfig{
			Host:           getEnv("MONGO_HOST", defaultMongoHost),
			Port:           getEnv("MONGO_PORT", defaultMongoPort),
			DBName:         getEnv("MONGO_DB", defaultMongoDB),
			User:           getEnv("MONGO_USER", defaultMongoUser),
			Pass:           getEnv("MONGO_PASS", defaultMongoPass),
			MaxRetries:     getEnvAsInt("MONGO_MAX_RETRIES", defaultMongoMaxRetries),
			RetryDelay:     getEnvAsDuration("MONGO_RETRY_DELAY_MS", defaultMongoRetryDelay),
			ConnectTimeout: getEnvAsDuration("MONGO_CONNECT_TIMEOUT_MS", defaultMongoConnectTimeout),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", defaultJWTSecret),
			Expiration:             getEnvAsDuration("JWT_EXPIRATION_MS", defaultJWTExpiration),
			RefreshTokenExpiration: getEnvAsDuration("REFRESH_TOKEN_EXPIRATION_MS", defaultRefreshTokenExpiration),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", defaultRedisHost),
			Port: getEnv("REDIS_PORT", defaultRedisPort),
			Pass: getEnv("REDIS_PASS", defaultRedisPass),
			DB:   getEnvAsInt("REDIS_DB", defaultRedisDB),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := getEnvAsInt(key, 0)
	if value == 0 {
		return defaultValue
	}

	return time.Duration(value) * time.Millisecond
}
