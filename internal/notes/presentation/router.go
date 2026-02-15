package presentation

import (
	"fmt"
	"log/slog"
	"net/http"

	"pinnado/internal/auth/infrastructure"
	"pinnado/pkg/nethttp"
)

type SetupRouterOptions struct {
	Mux         *http.ServeMux
	Prefix      string
	NoteService NoteService
	JWTService  interface {
		Validate(tokenString string) (*infrastructure.Claims, error)
	}
	Logger *slog.Logger
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewNoteHandler(options.NoteService)

	authMiddleware := nethttp.AuthMiddleware(options.JWTService)
	loggingMiddleware := makeLoggingMiddleware(options.Logger)

	createNoteHandler := addMiddlewares(
		handler.CreateNote,
		loggingMiddleware,
		authMiddleware,
	)

	listNotesHandler := addMiddlewares(
		handler.ListNotes,
		loggingMiddleware,
		authMiddleware,
	)

	options.Mux.Handle(
		fmt.Sprintf("POST %s/notes", options.Prefix),
		createNoteHandler,
	)

	options.Mux.Handle(
		fmt.Sprintf("GET %s/notes", options.Prefix),
		listNotesHandler,
	)
}

type middlewareFunc func(http.Handler) http.Handler

func addMiddlewares(handler http.HandlerFunc, middlewares ...middlewareFunc) http.Handler {
	result := http.Handler(handler)
	for i := len(middlewares) - 1; i >= 0; i-- {
		result = middlewares[i](result)
	}

	return result
}

func makeLoggingMiddleware(logger *slog.Logger) middlewareFunc {
	return func(handler http.Handler) http.Handler {
		if logger != nil {
			return nethttp.LoggingMiddleware(logger)(handler)
		}

		return handler
	}
}
