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

	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/refresh", options.Prefix),
		chain.ThenFunc(handler.Refresh),
	)

	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/logout", options.Prefix),
		chain.ThenFunc(handler.Logout),
	)

	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/forgot-password", options.Prefix),
		chain.ThenFunc(handler.ForgotPassword),
	)

	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/reset-password", options.Prefix),
		chain.ThenFunc(handler.ResetPassword),
	)

	options.Mux.Handle(
		fmt.Sprintf("GET %s/auth/reset-password/validate", options.Prefix),
		chain.ThenFunc(handler.ValidateResetToken),
	)
}
