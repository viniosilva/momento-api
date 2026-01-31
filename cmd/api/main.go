package main

import (
	"log"
	"net/http"

	"pinnado/internal/shared/application"
	"pinnado/internal/shared/presentation"
)

func main() {
	healthService := application.NewHealthService()

	mux := http.NewServeMux()
	presentation.SetupHealthRouter(mux, "/api", healthService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
