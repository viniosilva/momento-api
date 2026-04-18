package nethttp_timeout_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	nethttp_timeout "momento/pkg/nethttp/timeout"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware(t *testing.T) {
	t.Run("should work when handler completes before timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("completed"))
		})

		middleware := nethttp_timeout.TimeoutMiddleware(1 * time.Second)
		timeoutHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		timeoutHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "completed", rec.Body.String())
	})

	t.Run("should cancel context when timeout exceeds", func(t *testing.T) {
		var ctxErr error
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Millisecond)
			ctxErr = r.Context().Err()
		})

		middleware := nethttp_timeout.TimeoutMiddleware(1 * time.Millisecond)
		timeoutHandler := middleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		timeoutHandler.ServeHTTP(rec, req)
		assert.Error(t, ctxErr)
		assert.Equal(t, context.DeadlineExceeded, ctxErr)
	})
}
