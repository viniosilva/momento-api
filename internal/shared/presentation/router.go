package presentation

import (
	"fmt"
	"net/http"
	"path/filepath"
	"pinnado/internal/shared/application"

	httpSwagger "github.com/swaggo/http-swagger"
)

type SetupRouterOptions struct {
	Mux           *http.ServeMux
	Prefix        string
	HealthService *application.HealthService
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewHealthHandler(options.HealthService)

	serveSwaggerJSON(options)

	options.Mux.HandleFunc(fmt.Sprintf("GET %s/healthcheck", options.Prefix), handler.HealthCheck)
}

func serveSwaggerJSON(options SetupRouterOptions) {
	swaggerJSONPath := filepath.Join("docs", "swagger.json")

	options.Mux.HandleFunc("GET /docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerJSONPath)
	})
	options.Mux.HandleFunc("GET /docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))
	options.Mux.HandleFunc("GET /docs/*any", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))
}
