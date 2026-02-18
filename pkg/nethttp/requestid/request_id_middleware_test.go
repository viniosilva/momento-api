package nethttp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	nethttp_requestid "pinnado/pkg/nethttp/requestid"

	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("should generate request ID when not provided", func(t *testing.T) {
		var ctxFromHandler context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxFromHandler = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		nethttp_requestid.RequestIDMiddleware(handler).ServeHTTP(rec, req)

		requestID := nethttp_requestid.GetRequestID(ctxFromHandler)
		assert.NotEmpty(t, requestID)
		assert.NotEqual(t, "unknown", requestID)
		assert.Equal(t, requestID, rec.Header().Get("X-Request-ID"))
	})

	t.Run("should use provided request ID from header", func(t *testing.T) {
		var ctxFromHandler context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxFromHandler = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		providedID := "test-request-id-123"
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", providedID)
		rec := httptest.NewRecorder()

		nethttp_requestid.RequestIDMiddleware(handler).ServeHTTP(rec, req)

		requestID := nethttp_requestid.GetRequestID(ctxFromHandler)
		assert.Equal(t, providedID, requestID)
		assert.Equal(t, providedID, rec.Header().Get("X-Request-ID"))
	})
}

func TestGetRequestID(t *testing.T) {
	t.Run("should return request ID from context", func(t *testing.T) {
		expectedID := "test-id-456"

		ctx := context.WithValue(t.Context(), nethttp_requestid.RequestIDKey, expectedID)

		requestID := nethttp_requestid.GetRequestID(ctx)
		assert.Equal(t, expectedID, requestID)
	})

	t.Run("should return unknown when request ID not in context", func(t *testing.T) {
		requestID := nethttp_requestid.GetRequestID(t.Context())

		assert.Equal(t, "unknown", requestID)
	})
}
