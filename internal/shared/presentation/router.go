package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"pinnado/internal/shared/application"
	"pinnado/pkg/nethttp"

	httpSwagger "github.com/swaggo/http-swagger"
)

type SetupRouterOptions struct {
	Mux           *http.ServeMux
	Prefix        string
	HealthService *application.HealthService
	Logger        *slog.Logger
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewHealthHandler(options.HealthService)
	loggingMiddleware := makeLoggingMiddleware(options.Logger)

	serveSwaggerJSON(options)

	healthCheckHandler := addMiddleware(handler.HealthCheck, loggingMiddleware)
	options.Mux.Handle(fmt.Sprintf("GET %s/healthcheck", options.Prefix), healthCheckHandler)
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

type middlewareFunc func(http.Handler) http.Handler

func addMiddleware(handler http.HandlerFunc, middleware middlewareFunc) http.Handler {
	return middleware(handler)
}

func makeLoggingMiddleware(logger *slog.Logger) middlewareFunc {
	return func(handler http.Handler) http.Handler {
		if logger != nil {
			return nethttp.LoggingMiddleware(logger)(handler)
		}

		return handler
	}
}
