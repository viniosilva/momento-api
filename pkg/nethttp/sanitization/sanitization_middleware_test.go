package nethttp_sanitization_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	nethttp_sanitization "momento/pkg/nethttp/sanitization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizationMiddleware(t *testing.T) {
	t.Run("should remove XSS payload from string", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)

			content := body["content"].(string)
			assert.NotContains(t, content, "<script>")
			assert.NotContains(t, content, "alert")
			assert.Contains(t, content, "Hello")
			assert.Contains(t, content, "World")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}))

		payload := map[string]string{
			"content": "Hello <script>alert('xss')</script> World",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should remove HTML img onerror XSS", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)

			content := body["content"].(string)
			assert.NotContains(t, content, "img")
			assert.NotContains(t, content, "onerror")

			w.WriteHeader(http.StatusOK)
		}))

		payload := map[string]string{
			"content": "Click <img src=x onerror='alert(1)'>",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should reject string exceeding max size (5KB)", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		largeString := bytes.Repeat([]byte("a"), 6000)
		payload := map[string]any{
			"content": string(largeString),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "string too long")
	})

	t.Run("should reject body exceeding max size (1MB)", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		largePayload := make(map[string]string)
		largePayload["data"] = string(bytes.Repeat([]byte("x"), 1100000))

		body, _ := json.Marshal(largePayload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = int64(len(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "request body too large")
	})

	t.Run("should sanitize nested objects recursively", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)

			metadata := body["metadata"].(map[string]any)
			tags := metadata["tags"].([]any)

			assert.NotContains(t, tags[0], "<script>")
			assert.NotContains(t, tags[1], "img")

			w.WriteHeader(http.StatusOK)
		}))

		payload := map[string]any{
			"content": "Valid content",
			"metadata": map[string]any{
				"tags": []string{
					"security<script>alert('xss')</script>",
					"golang<img src=x onerror='alert()'>",
				},
			},
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bypass non-JSON requests", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		called := false
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader([]byte("text data")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.True(t, called, "handler should be called for non-JSON requests")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bypass GET requests", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		called := false
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/notes", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.True(t, called, "handler should be called for GET requests")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should reject invalid JSON", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader([]byte("{ invalid json }")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid JSON format")
	})

	t.Run("should accept valid data unchanged", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)

			assert.Equal(t, "Valid content", body["content"])
			assert.Equal(t, float64(42), body["count"])
			assert.Equal(t, true, body["active"])

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}))

		payload := map[string]any{
			"content": "Valid content",
			"count":   42,
			"active":  true,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should reject nesting deeper than 10 levels", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// 11 levels deep triggers the depth check (depth > 10)
		deepJSON := `{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":"deep"}}}}}}}}}}}`

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader([]byte(deepJSON)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "nesting too deep")
	})

	t.Run("should reject key longer than 100 chars", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		longKey := bytes.Repeat([]byte("k"), 101)
		payload := map[string]any{
			string(longKey): "value",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "key too long")
	})

	t.Run("should remove control characters from strings", func(t *testing.T) {
		middleware := nethttp_sanitization.SanitizationMiddleware()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)

			content := body["content"].(string)
			assert.NotContains(t, content, "\x01")
			assert.NotContains(t, content, "\x1f")
			assert.NotContains(t, content, "\x7f")
			assert.Contains(t, content, "Hello")
			assert.Contains(t, content, "\n")
			assert.Contains(t, content, "\t")

			w.WriteHeader(http.StatusOK)
		}))

		payload := map[string]any{
			"content": "Hello\x01World\x1f\x7f\nNewline\tTab",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
