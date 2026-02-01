package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

const (
	traceIDKey   = "trace_id"
	requestIDKey = "request_id"
)

func NewLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(newContextHandler(handler))

	return logger
}

type contextHandler struct {
	slog.Handler
}

func newContextHandler(h slog.Handler) *contextHandler {
	return &contextHandler{Handler: h}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := getValueFromContext(ctx, traceIDKey); traceID != "" {
		r.AddAttrs(slog.String(traceIDKey, traceID))
	}

	if requestID := getValueFromContext(ctx, requestIDKey); requestID != "" {
		r.AddAttrs(slog.String(requestIDKey, requestID))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newContextHandler(h.Handler.WithAttrs(attrs))
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return newContextHandler(h.Handler.WithGroup(name))
}

func getValueFromContext(ctx context.Context, key string) string {
	if val := ctx.Value(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}

	return ""
}
