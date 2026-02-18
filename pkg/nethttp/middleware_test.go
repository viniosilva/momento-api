package nethttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pinnado/pkg/nethttp"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	t.Run("should apply middlewares in correct order", func(t *testing.T) {
		var order []string

		middleware1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "m1-before")
				next.ServeHTTP(w, r)
				order = append(order, "m1-after")
			})
		}

		middleware2 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "m2-before")
				next.ServeHTTP(w, r)
				order = append(order, "m2-after")
			})
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "handler")
			w.WriteHeader(http.StatusOK)
		})

		chain := nethttp.NewChain(middleware1, middleware2)
		chainedHandler := chain.Then(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chainedHandler.ServeHTTP(rec, req)

		expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
		assert.Equal(t, expected, order)
	})

	t.Run("should apply middlewares with ThenFunc", func(t *testing.T) {
		called := false
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			})
		}

		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		chain := nethttp.NewChain(middleware)
		chainedHandler := chain.ThenFunc(handlerFunc)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chainedHandler.ServeHTTP(rec, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should apply middlewares with empty chain", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		chain := nethttp.NewChain()
		chainedHandler := chain.Then(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chainedHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})
}

func TestChain_AddMiddleware(t *testing.T) {
	t.Run("should append middleware and apply in order", func(t *testing.T) {
		var order []string

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "handler")
			w.WriteHeader(http.StatusOK)
		})

		extraMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "extra-before")
				next.ServeHTTP(w, r)
				order = append(order, "extra-after")
			})
		}

		chain := nethttp.NewChain()
		chain.AddMiddleware(extraMiddleware)
		chainedHandler := chain.Then(handler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chainedHandler.ServeHTTP(rec, req)

		assert.Equal(t, []string{"extra-before", "handler", "extra-after"}, order)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestNewDefaultChain(t *testing.T) {
	t.Run("should build chain with recovery, requestid and timeout", func(t *testing.T) {
		chain := nethttp.NewDefaultChain(nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chain.ThenFunc(handler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
		assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
	})

	t.Run("should recover from panic", func(t *testing.T) {
		chain := nethttp.NewDefaultChain(nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chain.ThenFunc(handler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
	})
}

func TestWithTimeout(t *testing.T) {
	t.Run("should use custom timeout when option provided", func(t *testing.T) {
		customTimeout := 5 * time.Second
		chain := nethttp.NewDefaultChain(nil, nethttp.WithTimeout(&customTimeout))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		chain.ThenFunc(handler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
