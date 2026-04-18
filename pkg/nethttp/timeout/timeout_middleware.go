package nethttp_timeout

import (
	"context"
	"net/http"
	nethttp_utils "momento/pkg/nethttp/utils"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) nethttp_utils.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
