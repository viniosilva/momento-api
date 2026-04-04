package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"pinnado/pkg/nethttp"
	auth "pinnado/pkg/nethttp/auth"
	logging "pinnado/pkg/nethttp/logging"
	sanitization "pinnado/pkg/nethttp/sanitization"
	"time"
)

type SetupRouterOptions struct {
	Mux         *http.ServeMux
	Prefix      string
	NoteService NoteService
	JWTService  JWTService
	Logger      *slog.Logger
	Timeout     *time.Duration
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewNoteHandler(options.NoteService)

	chain := nethttp.NewDefaultChain(options.Logger, nethttp.WithTimeout(options.Timeout))
	chain.AddMiddleware(sanitization.SanitizationMiddleware())
	chain.AddMiddleware(logging.LoggingMiddleware(options.Logger))
	chain.AddMiddleware(auth.AuthMiddleware(options.JWTService))

	options.Mux.Handle(
		fmt.Sprintf("POST %s/notes", options.Prefix),
		chain.ThenFunc(handler.CreateNote),
	)

	options.Mux.Handle(
		fmt.Sprintf("GET %s/notes", options.Prefix),
		chain.ThenFunc(handler.ListNotes),
	)

	options.Mux.Handle(
		fmt.Sprintf("GET %s/notes/{id}", options.Prefix),
		chain.ThenFunc(handler.GetUserNoteByID),
	)

	options.Mux.Handle(
		fmt.Sprintf("PATCH %s/notes/{id}", options.Prefix),
		chain.ThenFunc(handler.UpdateNote),
	)

	options.Mux.Handle(
		fmt.Sprintf("DELETE %s/notes/{id}", options.Prefix),
		chain.ThenFunc(handler.DeleteNote),
	)

	options.Mux.Handle(
		fmt.Sprintf("PATCH %s/notes/{id}/archive", options.Prefix),
		chain.ThenFunc(handler.ArchiveNote),
	)
}
