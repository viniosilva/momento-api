package main

import (
	"context"
	"log"
	"log/slog"
	authinfra "momento/internal/auth/infrastructure"
	"momento/internal/config"
	notesinfra "momento/internal/notes/infrastructure"
	"momento/pkg/logger"
	"momento/pkg/mongodb"
	"time"
)

const (
	shutdownTimeout = 10 * time.Second
)

func main() {
	slog.SetDefault(logger.NewLogger("info"))

	log.Println("loading configuration...")
	cfg := config.LoadConfig()

	log.Println("connecting to MongoDB...")
	ctx := context.Background()
	mongoClient, err := mongodb.NewMongoClient(ctx,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		cfg.Mongo.DBName,
		cfg.Mongo.User,
		cfg.Mongo.Pass,
		cfg.Mongo.MaxRetries,
		cfg.Mongo.RetryDelay,
		cfg.Mongo.ConnectTimeout)
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
	if err := authinfra.CreateIndexes(ctx, mongoClient.Database(cfg.Mongo.DBName)); err != nil {
		log.Fatalf("failed to create MongoDB indexes: %v", err)
	}
	if err := notesinfra.CreateIndexes(ctx, mongoClient.Database(cfg.Mongo.DBName)); err != nil {
		log.Fatalf("failed to create notes indexes: %v", err)
	}

	log.Println("migration completed successfully")
}
