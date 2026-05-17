package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Api               ApiConfig
	PG                PgConfig
	JWT               JWTConfig
	Redis             RedisConfig
	SMTP              SMTPConfig
	S3                S3Config
	ResetPassword     ResetPasswordConfig
	EmailVerification EmailVerificationConfig
}

type S3Config struct {
	Endpoint                   string
	Region                     string
	Bucket                     string
	AccessKey                  string
	SecretKey                  string
	UsePathStyle               bool
	UseSSL                     bool
	ImageDownloadURLExpiration time.Duration
}

type ApiConfig struct {
	Host string
	Port string
}

type PgConfig struct {
	DSN            string
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

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

type ResetPasswordConfig struct {
	TokenSize       int
	TokenExpiration time.Duration
	URLBase         string
}

type EmailVerificationConfig struct {
	TokenSize       int
	TokenExpiration time.Duration
	URLBase         string
}

const (
	defaultEnvPath                      = ".env"
	defaultApiHost                      = ""
	defaultApiPort                      = "8080"
	defaultPGDSN                        = "postgres://momento:momento@localhost:5432/momento?sslmode=disable"
	defaultPGMaxRetries                 = 3
	defaultPGRetryDelay                 = 2 * time.Second
	defaultPGConnectTimeout             = 10 * time.Second
	defaultJWTSecret                    = "your-secret-key-change-in-production"
	defaultJWTExpiration                = 12 * time.Hour
	defaultRefreshTokenExpiration       = 7 * 24 * time.Hour
	defaultRedisHost                    = "localhost"
	defaultRedisPort                    = "6379"
	defaultRedisPass                    = ""
	defaultRedisDB                      = 0
	defaultSMTPHost                     = "localhost"
	defaultSMTPPort                     = "1025"
	defaultSMTPUser                     = ""
	defaultSMTPPass                     = ""
	defaultSMTPFrom                     = "noreply@momento.com"
	defaultResetTokenSize               = 32
	defaultResetTokenExpiration         = 1 * time.Hour
	defaultResetURLBase                 = "http://http://momentonow.com/reset-password"
	defaultVerificationTokenSize        = 32
	defaultVerificationTokenExpiration  = 24 * time.Hour
	defaultVerificationURLBase          = "http://momentonow.com/sign-in"
	defaultS3Endpoint                   = "localhost:9000"
	defaultS3Region                     = "us-east-1"
	defaultS3Bucket                     = "momento"
	defaultS3AccessKey                  = "momento_admin"
	defaultS3SecretKey                  = "momento_admin"
	defaultS3UsePathStyle               = true
	defaultS3UseSSL                     = false
	defaultS3ImageDownloadURLExpiration = 15 * time.Minute
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
		PG: PgConfig{
			DSN:            getEnv("DATABASE_URL", defaultPGDSN),
			MaxRetries:     getEnvAsInt("PG_MAX_RETRIES", defaultPGMaxRetries),
			RetryDelay:     getEnvAsDuration("PG_RETRY_DELAY_MS", defaultPGRetryDelay),
			ConnectTimeout: getEnvAsDuration("PG_CONNECT_TIMEOUT_MS", defaultPGConnectTimeout),
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
		SMTP: SMTPConfig{
			Host: getEnv("SMTP_HOST", defaultSMTPHost),
			Port: getEnv("SMTP_PORT", defaultSMTPPort),
			User: getEnv("SMTP_USER", defaultSMTPUser),
			Pass: getEnv("SMTP_PASS", defaultSMTPPass),
			From: getEnv("SMTP_FROM", defaultSMTPFrom),
		},
		S3: S3Config{
			Endpoint:                   getEnv("S3_ENDPOINT", defaultS3Endpoint),
			Region:                     getEnv("S3_REGION", defaultS3Region),
			Bucket:                     getEnv("S3_BUCKET", defaultS3Bucket),
			AccessKey:                  getEnv("S3_ACCESS_KEY", defaultS3AccessKey),
			SecretKey:                  getEnv("S3_SECRET_KEY", defaultS3SecretKey),
			UsePathStyle:               getEnvAsBool("S3_USE_PATH_STYLE", defaultS3UsePathStyle),
			UseSSL:                     getEnvAsBool("S3_USE_SSL", defaultS3UseSSL),
			ImageDownloadURLExpiration: getEnvAsDuration("S3_IMAGE_DOWNLOAD_URL_EXPIRATION_MS", defaultS3ImageDownloadURLExpiration),
		},
		ResetPassword: ResetPasswordConfig{
			TokenSize:       getEnvAsInt("RESET_TOKEN_SIZE", defaultResetTokenSize),
			TokenExpiration: getEnvAsDuration("RESET_TOKEN_EXPIRATION_MS", defaultResetTokenExpiration),
			URLBase:         getEnv("RESET_URL_BASE", defaultResetURLBase),
		},
		EmailVerification: EmailVerificationConfig{
			TokenSize:       getEnvAsInt("VERIFICATION_TOKEN_SIZE", defaultVerificationTokenSize),
			TokenExpiration: getEnvAsDuration("VERIFICATION_TOKEN_EXPIRATION_MS", defaultVerificationTokenExpiration),
			URLBase:         getEnv("VERIFICATION_URL_BASE", defaultVerificationURLBase),
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

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
