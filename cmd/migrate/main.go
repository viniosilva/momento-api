package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	authinfra "pinnado/internal/auth/infrastructure"
	notesinfra "pinnado/internal/notes/infrastructure"
	"pinnado/internal/shared/infrastructure"
	"pinnado/pkg/logger"
	"pinnado/pkg/mongodb"
)

const (
	shutdownTimeout = 10 * time.Second
)

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
	if err := notesinfra.CreateIndexes(ctx, mongoClient.Database(config.Mongo.DBName)); err != nil {
		log.Fatalf("failed to create notes indexes: %v", err)
	}

	log.Println("migration completed successfully")
}
