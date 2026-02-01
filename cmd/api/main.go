package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pinnado/docs"
	"pinnado/internal/shared/application"
	"pinnado/internal/shared/infrastructure"
	"pinnado/internal/shared/presentation"
	"pinnado/pkg/mongodb"
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

	addr := fmt.Sprintf("%s:%s", config.Api.Host, config.Api.Port)
	docs.SwaggerInfo.Host = addr

	mux := http.NewServeMux()
	presentation.SetupRouter(presentation.SetupRouterOptions{
		Mux:           mux,
		Prefix:        "/api",
		HealthService: healthService,
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("server starting on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
