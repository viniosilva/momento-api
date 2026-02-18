package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"pinnado/pkg/nethttp"
	logging "pinnado/pkg/nethttp/logging"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

type SetupRouterOptions struct {
	Mux           *http.ServeMux
	Prefix        string
	HealthService HealthService
	Logger        *slog.Logger
	Timeout       *time.Duration
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewHealthHandler(options.HealthService)

	chain := nethttp.NewDefaultChain(options.Logger, nethttp.WithTimeout(options.Timeout))
	chain.AddMiddleware(logging.LoggingMiddleware(options.Logger))

	healthCheckHandler := chain.ThenFunc(handler.HealthCheck)

	serveSwaggerJSON(options)
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
