package main

import (
	"context"
	"fmt"
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
	notesapp "pinnado/internal/notes/application"
	notesdomain "pinnado/internal/notes/domain"
	notesinfra "pinnado/internal/notes/infrastructure"
	notespres "pinnado/internal/notes/presentation"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/infrastructure"
	"pinnado/internal/shared/presentation"
	pkglogger "pinnado/pkg/logger"
	"pinnado/pkg/mongodb"
)

const (
	shutdownTimeout = 10 * time.Second
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
	config := infrastructure.LoadConfig()

	logger.Info("connecting to MongoDB")
	ctx := context.Background()
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
		logger.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer func() {
		logger.Info("disconnecting from MongoDB")
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("error disconnecting from MongoDB", "error", err)
		}
	}()

	logger.Info("initializing services")
	healthService := application.NewHealthService(mongoClient)

	db := mongoClient.Database(config.Mongo.DBName)
	userCollection := db.Collection(authdomain.UsersCollectionName)
	userRepository := authinfra.NewUserRepository(userCollection)
	jwtService := authinfra.NewJWTService(config.JWT.Secret, config.JWT.Expiration)
	authService := authapp.NewAuthService(userRepository, jwtService)

	noteCollection := db.Collection(notesdomain.NotesCollectionName)
	noteRepository := notesinfra.NewNoteRepository(noteCollection)
	noteService := notesapp.NewNoteService(noteRepository)

	addr := fmt.Sprintf("%s:%s", config.Api.Host, config.Api.Port)
	docs.SwaggerInfo.Host = addr

	mux := http.NewServeMux()

	presentation.SetupRouter(presentation.SetupRouterOptions{
		Mux:           mux,
		Prefix:        "/api",
		HealthService: healthService,
		Logger:        logger,
	})

	authpres.SetupRouter(authpres.SetupRouterOptions{
		Mux:         mux,
		Prefix:      "/api",
		AuthService: authService,
		Logger:      logger,
	})

	notespres.SetupRouter(notespres.SetupRouterOptions{
		Mux:         mux,
		Prefix:      "/api",
		NoteService: noteService,
		JWTService:  jwtService,
		Logger:      logger,
	})

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
