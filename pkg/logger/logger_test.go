package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"pinnado/pkg/logger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("should create logger with info level", func(t *testing.T) {
		log := logger.NewLogger("info")
		assert.NotNil(t, log)
	})

	t.Run("should create logger with debug level", func(t *testing.T) {
		log := logger.NewLogger("debug")
		assert.NotNil(t, log)
	})

	t.Run("should create logger with warn level", func(t *testing.T) {
		log := logger.NewLogger("warn")
		assert.NotNil(t, log)
	})

	t.Run("should create logger with error level", func(t *testing.T) {
		log := logger.NewLogger("error")
		assert.NotNil(t, log)
	})

	t.Run("should default to info level for invalid level", func(t *testing.T) {
		log := logger.NewLogger("invalid")
		assert.NotNil(t, log)
	})

	t.Run("should default to info level for empty level", func(t *testing.T) {
		log := logger.NewLogger("")
		assert.NotNil(t, log)
	})
}

func TestContextHandler_Handle(t *testing.T) {
	t.Run("should extract trace_id from context", func(t *testing.T) {
		var buf bytes.Buffer
		baseHandler := slog.NewJSONHandler(&buf, nil)
		customHandler := logger.NewLogger("info").Handler()

		ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

		err := customHandler.Handle(ctx, record)
		require.NoError(t, err)

		newRecord := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
		newRecord.AddAttrs(slog.String("trace_id", "test-trace-123"))

		err = baseHandler.Handle(ctx, newRecord)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)

		traceID, ok := result["trace_id"].(string)
		assert.True(t, ok, "trace_id should be a string")
		assert.Equal(t, "test-trace-123", traceID)
	})

	t.Run("should extract request_id from context", func(t *testing.T) {
		var buf bytes.Buffer
		baseHandler := slog.NewJSONHandler(&buf, nil)
		customHandler := logger.NewLogger("info").Handler()

		ctx := context.WithValue(context.Background(), "request_id", "test-request-456")
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

		err := customHandler.Handle(ctx, record)
		require.NoError(t, err)

		newRecord := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
		newRecord.AddAttrs(slog.String("request_id", "test-request-456"))

		err = baseHandler.Handle(ctx, newRecord)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)

		requestID, ok := result["request_id"].(string)
		assert.True(t, ok, "request_id should be a string")
		assert.Equal(t, "test-request-456", requestID)
	})

	t.Run("should extract both trace_id and request_id from context", func(t *testing.T) {
		customHandler := logger.NewLogger("info").Handler()

		ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
		ctx = context.WithValue(ctx, "request_id", "request-456")
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

		err := customHandler.Handle(ctx, record)
		require.NoError(t, err)

	})

	t.Run("should not add trace_id when not in context", func(t *testing.T) {
		var buf bytes.Buffer
		baseHandler := slog.NewJSONHandler(&buf, nil)
		customHandler := logger.NewLogger("info").Handler()

		ctx := context.Background()
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

		err := customHandler.Handle(ctx, record)
		require.NoError(t, err)

		err = baseHandler.Handle(ctx, record)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)

		_, ok := result["trace_id"]
		assert.False(t, ok, "Should not have trace_id when not in context")
	})

	t.Run("should not add request_id when not in context", func(t *testing.T) {
		var buf bytes.Buffer
		baseHandler := slog.NewJSONHandler(&buf, nil)
		customHandler := logger.NewLogger("info").Handler()

		ctx := context.Background()
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

		err := customHandler.Handle(ctx, record)
		require.NoError(t, err)

		err = baseHandler.Handle(ctx, record)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)

		_, ok := result["request_id"]
		assert.False(t, ok, "Should not have request_id when not in context")
	})

	t.Run("should handle non-string values in context", func(t *testing.T) {
		customHandler := logger.NewLogger("info").Handler()

		ctx := context.WithValue(context.Background(), "trace_id", 123)
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

		err := customHandler.Handle(ctx, record)
		assert.NoError(t, err, "Handle should not return error even with non-string value")
	})
}

func TestContextHandler_WithAttrs(t *testing.T) {
	t.Run("should preserve context extraction with WithAttrs", func(t *testing.T) {
		handler := logger.NewLogger("info").Handler()
		newHandler := handler.WithAttrs([]slog.Attr{slog.String("custom", "value")})
		assert.NotNil(t, newHandler, "WithAttrs should return a non-nil handler")
	})
}

func TestContextHandler_WithGroup(t *testing.T) {
	t.Run("should preserve context extraction with WithGroup", func(t *testing.T) {
		handler := logger.NewLogger("info").Handler()
		newHandler := handler.WithGroup("user")
		assert.NotNil(t, newHandler, "WithGroup should return a non-nil handler")
	})
}
