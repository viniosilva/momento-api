package nethttp

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	t.Run("should return 200 ok response", func(t *testing.T) {
		body := map[string]string{"status": "ok"}
		cb := func(w http.ResponseWriter, r *http.Request) error {
			JSON(w, http.StatusOK, body)
			return nil
		}

		resp, err := Request(t.Context(), http.MethodGet, "/test", body, cb)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	t.Run("should return 500 error response", func(t *testing.T) {
		body := map[string]string{"error": "internal server error"}
		expectedErr := errors.New("internal server error")
		cb := func(w http.ResponseWriter, r *http.Request) error {
			JSON(w, http.StatusInternalServerError, body)
			return expectedErr
		}

		resp, err := Request(t.Context(), http.MethodGet, "/test", body, cb)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
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
		cb := func(w http.ResponseWriter, r *http.Request) error {
			return nil
		}

		resp, err := Request(t.Context(), http.MethodGet, "/test", body, cb)

		assert.Error(t, err)
		assert.Nil(t, resp)
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

		cb := func(w http.ResponseWriter, r *http.Request) error {
			JSON(w, http.StatusOK, responseBody)
			return nil
		}

		type ResponseType struct {
			Status string `json:"status"`
		}

		resp, got, err := RequestWithResponse[map[string]string, ResponseType](t.Context(), http.MethodGet, "/test", requestBody, cb)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.NotNil(t, got)
		assert.Equal(t, "ok", got.Status)
	})

	t.Run("should return error when json decode fails", func(t *testing.T) {
		requestBody := map[string]string{"action": "check"}
		cb := func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
			return nil
		}

		type ResponseType struct {
			Status string `json:"status"`
		}

		resp, got, err := RequestWithResponse[map[string]string, ResponseType](t.Context(), http.MethodGet, "/test", requestBody, cb)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Nil(t, got)
	})

	t.Run("should return error when Request returns error", func(t *testing.T) {
		requestBody := map[string]string{"action": "check"}
		expectedErr := errors.New("callback error")
		cb := func(w http.ResponseWriter, r *http.Request) error {
			return expectedErr
		}

		type ResponseType struct {
			Status string `json:"status"`
		}

		resp, got, err := RequestWithResponse[map[string]string, ResponseType](t.Context(), http.MethodGet, "/test", requestBody, cb)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)
		assert.Nil(t, got)
	})
}
