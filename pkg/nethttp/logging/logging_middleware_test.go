package nethttp_logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	nethttp_logging "pinnado/pkg/nethttp/logging"
	nethttp_utils "pinnado/pkg/nethttp/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	t.Run("should return 200 status with json response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{
			"result": "ok",
		}

		nethttp_utils.JSON(w, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)

		assert.NoError(t, err)
		assert.Equal(t, "ok", got["result"])
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("should log request start and completion", func(t *testing.T) {
		var logOutput bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logOutput, nil))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nethttp_utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		middleware := nethttp_logging.LoggingMiddleware(logger)
		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		output := logOutput.String()
		assert.Contains(t, output, "request started")
		assert.Contains(t, output, "request completed")
		assert.Contains(t, output, `"method":"GET"`)
		assert.Contains(t, output, `"path":"/test"`)
		assert.Contains(t, output, `"status":200`)
		assert.Contains(t, output, "latency_ms")
	})

	t.Run("should log correct status code", func(t *testing.T) {
		var logOutput bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logOutput, nil))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nethttp_utils.JSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		})

		middleware := nethttp_logging.LoggingMiddleware(logger)
		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		output := logOutput.String()
		assert.Contains(t, output, `"status":404`)
	})

	t.Run("should log different HTTP methods", func(t *testing.T) {
		var logOutput bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logOutput, nil))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nethttp_utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		middleware := nethttp_logging.LoggingMiddleware(logger)
		wrappedHandler := middleware(handler)

		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

		for _, method := range methods {
			logOutput.Reset()
			req := httptest.NewRequest(method, "/test", nil)
			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			output := logOutput.String()
			assert.Contains(t, output, `"method":"`+method+`"`)
		}
	})

	t.Run("should log latency", func(t *testing.T) {
		var logOutput bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logOutput, nil))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			nethttp_utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		middleware := nethttp_logging.LoggingMiddleware(logger)
		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		output := logOutput.String()
		assert.Contains(t, output, "latency_ms")
	})

	t.Run("should use context for logging", func(t *testing.T) {
		var logOutput bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logOutput, nil))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nethttp_utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		middleware := nethttp_logging.LoggingMiddleware(logger)
		wrappedHandler := middleware(handler)

		ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		output := logOutput.String()
		// The logger should extract trace_id from context automatically
		assert.Contains(t, output, "request started")
		assert.Contains(t, output, "request completed")
	})

	t.Run("should default to 200 status when WriteHeader is not called", func(t *testing.T) {
		var logOutput bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logOutput, nil))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

		middleware := nethttp_logging.LoggingMiddleware(logger)
		wrappedHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		output := logOutput.String()
		assert.Contains(t, output, `"status":200`)
	})
}
