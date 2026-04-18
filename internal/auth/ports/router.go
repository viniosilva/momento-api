package ports

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"momento/pkg/nethttp"
	logging "momento/pkg/nethttp/logging"
)

type SetupRouterOptions struct {
	Mux         *http.ServeMux
	Prefix      string
	AuthService AuthService
	Logger      *slog.Logger
	Timeout     *time.Duration
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewAuthHandler(options.AuthService)

	chain := nethttp.NewDefaultChain(options.Logger, nethttp.WithTimeout(options.Timeout))
	chain.AddMiddleware(logging.LoggingMiddleware(options.Logger))

	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/register", options.Prefix),
		chain.ThenFunc(handler.Register),
	)

	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/login", options.Prefix),
		chain.ThenFunc(handler.Login),
	)
}
