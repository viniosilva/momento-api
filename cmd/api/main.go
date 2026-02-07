package main

import (
	"context"
	"fmt"
	"log"
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
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/infrastructure"
	"pinnado/internal/shared/presentation"
	"pinnado/pkg/logger"
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
// @BasePath /api
func main() {
	slog.SetDefault(logger.NewLogger("info"))

	log.Println("loading configuration...")
	config := infrastructure.LoadConfig()

	log.Println("connecting to MongoDB...")
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
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		log.Println("disconnecting from MongoDB...")
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("error disconnecting from MongoDB: %v", err)
		}
	}()

	log.Println("creating MongoDB indexes...")
	if err := authinfra.CreateIndexes(ctx, mongoClient.Database(config.Mongo.DBName)); err != nil {
		log.Fatalf("failed to create MongoDB indexes: %v", err)
	}

	log.Println("initializing services...")
	healthService := application.NewHealthService(mongoClient)

	db := mongoClient.Database(config.Mongo.DBName)
	userCollection := db.Collection(authdomain.UsersCollectionName)
	userRepository := authinfra.NewUserRepository(userCollection)
	authService := authapp.NewAuthService(userRepository)

	addr := fmt.Sprintf("%s:%s", config.Api.Host, config.Api.Port)
	docs.SwaggerInfo.Host = addr

	appLogger := logger.NewLogger("info")
	mux := http.NewServeMux()

	presentation.SetupRouter(presentation.SetupRouterOptions{
		Mux:           mux,
		Prefix:        "/api",
		HealthService: healthService,
		Logger:        appLogger,
	})

	authpres.SetupRouter(authpres.SetupRouterOptions{
		Mux:         mux,
		Prefix:      "/api",
		AuthService: authService,
		Logger:      appLogger,
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited gracefully")
}
