package nethttp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // default status code
			}

			requestID := uuid.New().String()
			ctx := setContextRequestID(r.Context(), requestID)

			logger.InfoContext(ctx,
				"request started",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(rw, r)

			latency := time.Since(start)

			logger.InfoContext(ctx,
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"latency_ms", latency.Milliseconds(),
			)
		})
	}
}

func setContextRequestID(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, "request_id", requestID)

	return ctx
}
