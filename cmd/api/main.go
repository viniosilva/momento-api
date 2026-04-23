package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"momento/docs"
	authadapters "momento/internal/auth/adapters"
	authapp "momento/internal/auth/app"
	authports "momento/internal/auth/ports"
	"momento/internal/config"
	eventsadapters "momento/internal/events/adapters"
	eventsapp "momento/internal/events/app"
	eventsports "momento/internal/events/ports"
	sharedapp "momento/internal/shared/app"
	sharedports "momento/internal/shared/ports"
	pkglogger "momento/pkg/logger"
	"momento/pkg/mongodb"
	pkgredis "momento/pkg/redis"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	shutdownTimeout = 10 * time.Second
	apiPrefixPath   = "/api"
)

// @title Momento API
// @version 1.0
// @description API documentation for Momento application
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@momento.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @BasePath /api
func main() {
	logger := pkglogger.NewLogger("info")

	logger.Info("loading configuration")
	cfg := config.LoadConfig()

	di, err := setupDependencies(context.Background(), cfg, logger)
	if err != nil {
		logger.Error("failed to setup dependencies", "error", err)
		os.Exit(1)
	}
	defer func() {
		logger.Info("disconnecting from MongoDB")
		if err := di.MongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("error disconnecting from MongoDB", "error", err)
		}
	}()

	addr := fmt.Sprintf("%s:%s", cfg.Api.Host, cfg.Api.Port)
	docs.SwaggerInfo.Host = addr

	mux := http.NewServeMux()
	setupRoutes(mux, di, logger)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		logger.Info("server starting", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("shutdown signal received", "signal", sig)

	shutdownStart := time.Now()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	logger.Info("shutting down server gracefully", "timeout", shutdownTimeout.String())

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err, "elapsed", time.Since(shutdownStart))
		os.Exit(1)
	}

	logger.Info("server exited gracefully", "elapsed", time.Since(shutdownStart))
}

type jwtSvc interface {
	authapp.JWTService
	eventsports.JWTService
}

type Dependencies struct {
	MongoClient   *mongo.Client
	HealthService sharedports.HealthService
	JwtService    jwtSvc
	AuthService   authports.AuthService
	EventService  eventsports.EventService
}

func setupDependencies(ctx context.Context, config config.Config, logger *slog.Logger) (*Dependencies, error) {
	logger.Info("connecting to MongoDB")
	mongoClient, err := mongodb.NewMongoClient(ctx,
		config.Mongo.Host,
		config.Mongo.Port,
		config.Mongo.DBName,
		config.Mongo.User,
		config.Mongo.Pass,
		config.Mongo.MaxRetries,
		config.Mongo.RetryDelay,
		config.Mongo.ConnectTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	logger.Info("connecting to Redis")
	redisClient, err := pkgredis.NewRedisClient(
		ctx,
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Pass,
		config.Redis.DB,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("initializing services")

	db := mongoClient.Database(config.Mongo.DBName)
	userRepository := authadapters.NewUserRepository(db)
	eventRepository := eventsadapters.NewEventRepository(db)
	secureTokenRepository := authadapters.NewSecureTokenService(redisClient, config.JWT.RefreshTokenExpiration)

	jwtService := authadapters.NewJWTService(config.JWT.Secret, config.JWT.Expiration)

	return &Dependencies{
		MongoClient:   mongoClient,
		HealthService: sharedapp.NewHealthService(mongoClient),
		JwtService:    jwtService,
		AuthService:   authapp.NewAuthService(userRepository, jwtService, secureTokenRepository),
		EventService:  eventsapp.NewEventService(eventRepository),
	}, nil
}

func setupRoutes(mux *http.ServeMux, di *Dependencies, logger *slog.Logger) {
	sharedports.SetupRouter(sharedports.SetupRouterOptions{
		Mux:           mux,
		Prefix:        apiPrefixPath,
		HealthService: di.HealthService,
		Logger:        logger,
	})

	authports.SetupRouter(authports.SetupRouterOptions{
		Mux:         mux,
		Prefix:      apiPrefixPath,
		AuthService: di.AuthService,
		Logger:      logger,
	})

	eventsports.SetupRouter(eventsports.SetupRouterOptions{
		Mux:          mux,
		Prefix:       apiPrefixPath,
		EventService: di.EventService,
		JWTService:   di.JwtService,
		Logger:       logger,
	})
}