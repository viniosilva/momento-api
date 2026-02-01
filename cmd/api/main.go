package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"pinnado/internal/shared/application"
	"pinnado/internal/shared/infrastructure"
	"pinnado/internal/shared/presentation"
	"pinnado/pkg/mongodb"
)

func main() {
	config := infrastructure.LoadConfig()

	mongoClient, err := mongodb.NewMongoClient(context.Background(),
		config.Mongo.Host,
		config.Mongo.Port,
		config.Mongo.DBName,
		config.Mongo.User,
		config.Mongo.Pass,
		config.Mongo.MaxRetries,
		config.Mongo.RetryDelay,
		config.Mongo.ConnectTimeout)
	if err != nil {
		log.Fatal(err)
	}

	healthService := application.NewHealthService(mongoClient)

	mux := http.NewServeMux()
	presentation.SetupHealthRouter(mux, "/api", healthService)

	addr := fmt.Sprintf("%s:%s", config.Api.Host, config.Api.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("server starting on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
