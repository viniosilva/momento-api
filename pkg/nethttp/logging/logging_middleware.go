package nethttp_logging

import (
	"net/http"
	nethttp_port "momento/pkg/nethttp/port"
	nethttp_utils "momento/pkg/nethttp/utils"
	"time"
)

func LoggingMiddleware(logger nethttp_port.LoggerSlog) nethttp_utils.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := nethttp_utils.NewResponseWriter(w)
			logger.InfoContext(r.Context(),
				"request started",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(rw, r)

			latency := time.Since(start)

			logger.InfoContext(r.Context(),
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.GetStatusCode(),
				"latency_ms", latency.Milliseconds(),
			)
		})
	}
}
