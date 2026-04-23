package ports

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"momento/pkg/nethttp"
	auth "momento/pkg/nethttp/auth"
	logging "momento/pkg/nethttp/logging"
	sanitization "momento/pkg/nethttp/sanitization"
)

type SetupRouterOptions struct {
	Mux         *http.ServeMux
	Prefix      string
	EventService EventService
	JWTService  JWTService
	Logger      *slog.Logger
	Timeout     *time.Duration
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewEventHandler(options.EventService)

	chain := nethttp.NewDefaultChain(options.Logger, nethttp.WithTimeout(options.Timeout))
	chain.AddMiddleware(logging.LoggingMiddleware(options.Logger))
	chain.AddMiddleware(sanitization.SanitizationMiddleware())
	chain.AddMiddleware(auth.AuthMiddleware(options.JWTService))

	options.Mux.Handle(
		fmt.Sprintf("POST %s/events", options.Prefix),
		chain.ThenFunc(handler.CreateEvent),
	)

	options.Mux.Handle(
		fmt.Sprintf("GET %s/events", options.Prefix),
		chain.ThenFunc(handler.ListEvents),
	)

	options.Mux.Handle(
		fmt.Sprintf("GET %s/events/{id}", options.Prefix),
		chain.ThenFunc(handler.GetUserEventByID),
	)

	options.Mux.Handle(
		fmt.Sprintf("PATCH %s/events/{id}", options.Prefix),
		chain.ThenFunc(handler.UpdateEvent),
	)

	options.Mux.Handle(
		fmt.Sprintf("PATCH %s/events/{id}/archive", options.Prefix),
		chain.ThenFunc(handler.ArchiveEvent),
	)

	options.Mux.Handle(
		fmt.Sprintf("PATCH %s/events/{id}/restore", options.Prefix),
		chain.ThenFunc(handler.RestoreEvent),
	)

	options.Mux.Handle(
		fmt.Sprintf("DELETE %s/events/{id}", options.Prefix),
		chain.ThenFunc(handler.DeleteEvent),
	)
}