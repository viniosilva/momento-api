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

	"pinnado/docs"
	authapp "pinnado/internal/auth/application"
	authdomain "pinnado/internal/auth/domain"
	authinfra "pinnado/internal/auth/infrastructure"
	authpres "pinnado/internal/auth/presentation"
	"pinnado/internal/config"
	notesapp "pinnado/internal/notes/application"
	notesdomain "pinnado/internal/notes/domain"
	notesinfra "pinnado/internal/notes/infrastructure"
	notespres "pinnado/internal/notes/presentation"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/presentation"
	pkglogger "pinnado/pkg/logger"
	"pinnado/pkg/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	shutdownTimeout = 10 * time.Second
	apiPrefixPath   = "/api"
)

// @title Pinnado API
// @version 1.0
// @description API documentation for Pinnado application
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@pinnado.com
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

type Dependencies struct {
	MongoClient   *mongo.Client
	HealthService *application.HealthService
	JwtService    *authinfra.JwtService
	AuthService   *authapp.AuthService
	NoteService   *notesapp.NoteService
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

	logger.Info("initializing services")

	db := mongoClient.Database(config.Mongo.DBName)
	userRepository := authinfra.NewUserRepository(db.Collection(authdomain.UsersCollectionName))
	noteRepository := notesinfra.NewNoteRepository(db.Collection(notesdomain.NotesCollectionName))

	jwtService := authinfra.NewJWTService(config.JWT.Secret, config.JWT.Expiration)

	return &Dependencies{
		MongoClient:   mongoClient,
		HealthService: application.NewHealthService(mongoClient),
		JwtService:    jwtService,
		AuthService:   authapp.NewAuthService(userRepository, jwtService),
		NoteService:   notesapp.NewNoteService(noteRepository),
	}, nil
}

func setupRoutes(mux *http.ServeMux, di *Dependencies, logger *slog.Logger) {
	presentation.SetupRouter(presentation.SetupRouterOptions{
		Mux:           mux,
		Prefix:        apiPrefixPath,
		HealthService: di.HealthService,
		Logger:        logger,
	})

	authpres.SetupRouter(authpres.SetupRouterOptions{
		Mux:         mux,
		Prefix:      apiPrefixPath,
		AuthService: di.AuthService,
		Logger:      logger,
	})

	notespres.SetupRouter(notespres.SetupRouterOptions{
		Mux:         mux,
		Prefix:      apiPrefixPath,
		NoteService: di.NoteService,
		JWTService:  di.JwtService,
		Logger:      logger,
	})
}
