package nethttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	t.Run("should return 200 ok response", func(t *testing.T) {
		body := map[string]string{"status": "ok"}
		cb := func(w http.ResponseWriter, r *http.Request) {
			JSON(w, http.StatusOK, body)
		}

		resp, err := Request(t.Context(), http.MethodGet, "/test", body, cb)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	t.Run("should return 500 error response", func(t *testing.T) {
		body := map[string]string{"error": "internal server error"}
		cb := func(w http.ResponseWriter, r *http.Request) {
			JSON(w, http.StatusInternalServerError, body)
		}

		resp, err := Request(t.Context(), http.MethodGet, "/test", body, cb)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	t.Run("should return error when json marshal fails", func(t *testing.T) {
		type invalidBody struct {
			Channel chan int
		}
		body := invalidBody{
			Channel: make(chan int),
		}
		cb := func(w http.ResponseWriter, r *http.Request) {
		}

		_, err := Request(t.Context(), http.MethodGet, "/test", body, cb)

		assert.ErrorContains(t, err, "json: unsupported type: chan int")
	})
}

func TestRequestWithResponse(t *testing.T) {
	t.Run("should return 200 ok response with decoded body", func(t *testing.T) {
		requestBody := map[string]string{"action": "check"}
		responseBody := struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		}

		cb := func(w http.ResponseWriter, r *http.Request) {
			JSON(w, http.StatusOK, responseBody)
		}

		type ResponseType struct {
			Status string `json:"status"`
		}

		resp, got, err := RequestWithResponse[map[string]string, ResponseType](t.Context(), http.MethodGet, "/test", requestBody, cb)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, "ok", got.Status)
	})

	t.Run("should return error when Request fails", func(t *testing.T) {
		type invalidBody struct {
			Channel chan int
		}
		body := invalidBody{
			Channel: make(chan int),
		}
		cb := func(w http.ResponseWriter, r *http.Request) {}

		type ResponseType struct {
			Status string `json:"status"`
		}

		_, _, err := RequestWithResponse[invalidBody, ResponseType](t.Context(), http.MethodGet, "/test", body, cb)

		assert.ErrorContains(t, err, "json: unsupported type: chan int")
	})

	t.Run("should return error when json decode fails", func(t *testing.T) {
		requestBody := map[string]string{"action": "check"}
		cb := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}

		type ResponseType struct {
			Status string `json:"status"`
		}

		_, _, err := RequestWithResponse[map[string]string, ResponseType](t.Context(), http.MethodGet, "/test", requestBody, cb)

		assert.ErrorContains(t, err, "invalid character")
	})
}
