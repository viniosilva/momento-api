package nethttp_recovery

import (
	"net/http"
	nethttp_port "momento/pkg/nethttp/port"
	nethttp_utils "momento/pkg/nethttp/utils"
)

func RecoveryMiddleware(logger nethttp_port.LoggerSlog) nethttp_utils.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if logger != nil {
						logger.ErrorContext(r.Context(), "panic recovered",
							"error", err,
							"path", r.URL.Path,
							"method", r.Method)
					}

					nethttp_utils.JSON(w, http.StatusInternalServerError,
						map[string]string{"message": "internal server error"})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
