package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionURI       = "mongodb://%s:%s@%s:%s/%s?authSource=admin"
	connectionURINoAuth = "mongodb://%s:%s/%s"
)

func NewMongoClient(ctx context.Context, host, port, dbName, user, pass string,
	maxRetries int, retryDelay, connectTimeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	var client *mongo.Client
	var err error

	uri := buildConnectionURI(host, port, dbName, user, pass)
	clientOptions := options.Client().ApplyURI(uri)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			if err = client.Ping(ctx, nil); err == nil {
				return client, nil
			}
			client.Disconnect(ctx)
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to MongoDB after %d attempts: %w", maxRetries, err)
}

func buildConnectionURI(host, port, dbName, user, pass string) string {
	if user != "" && pass != "" {
		return fmt.Sprintf(connectionURI, user, pass, host, port, dbName)
	}

	return fmt.Sprintf(connectionURINoAuth, host, port, dbName)
}
