package presentation

import (
	"fmt"
	"log/slog"
	"net/http"

	"pinnado/pkg/nethttp"
)

type SetupRouterOptions struct {
	Mux         *http.ServeMux
	Prefix      string
	AuthService AuthService
	Logger      *slog.Logger
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewAuthHandler(options.AuthService)
	loggingMiddleware := makeLoggingMiddleware(options.Logger)

	registerHandler := addMiddleware(handler.Register, loggingMiddleware)
	options.Mux.Handle(fmt.Sprintf("POST %s/auth/register", options.Prefix), registerHandler)
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
