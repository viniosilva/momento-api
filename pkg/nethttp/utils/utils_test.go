package nethttp_utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	nethttp_utils "momento/pkg/nethttp/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCode(t *testing.T) {
	t.Run("should set status code", func(t *testing.T) {
		rec := httptest.NewRecorder()
		nethttp_utils.StatusCode(rec, http.StatusOK)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestJSON(t *testing.T) {
	t.Run("should set Content-Type application/json", func(t *testing.T) {
		rec := httptest.NewRecorder()
		nethttp_utils.JSON(rec, http.StatusOK, map[string]string{"key": "value"})

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("should set status code", func(t *testing.T) {
		rec := httptest.NewRecorder()
		nethttp_utils.JSON(rec, http.StatusCreated, nil)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("should encode body as JSON", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body := map[string]any{"id": "123", "name": "test"}

		nethttp_utils.JSON(rec, http.StatusOK, body)

		assert.Equal(t, http.StatusOK, rec.Code)
		var decoded map[string]any
		err := json.NewDecoder(rec.Body).Decode(&decoded)
		require.NoError(t, err)
		assert.Equal(t, "123", decoded["id"])
		assert.Equal(t, "test", decoded["name"])
	})

	t.Run("should encode slice as JSON", func(t *testing.T) {
		rec := httptest.NewRecorder()
		body := []string{"a", "b", "c"}

		nethttp_utils.JSON(rec, http.StatusOK, body)

		var decoded []string
		err := json.NewDecoder(rec.Body).Decode(&decoded)
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, decoded)
	})
}

func TestNewResponseWriter(t *testing.T) {
	t.Run("should wrap ResponseWriter with default status 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := nethttp_utils.NewResponseWriter(rec)

		assert.Equal(t, http.StatusOK, rw.GetStatusCode())
	})
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	t.Run("should set status code and delegate to underlying ResponseWriter", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := nethttp_utils.NewResponseWriter(rec)

		rw.WriteHeader(http.StatusNotFound)

		assert.Equal(t, http.StatusNotFound, rw.GetStatusCode())
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should not panic on multiple WriteHeader calls", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := nethttp_utils.NewResponseWriter(rec)

		rw.WriteHeader(http.StatusCreated)
		rw.WriteHeader(http.StatusOK)

		// Underlying ResponseWriter ignores second call; wrapper may keep last
		assert.True(t, rw.GetStatusCode() == http.StatusCreated || rw.GetStatusCode() == http.StatusOK)
	})
}

func TestResponseWriter_GetStatusCode(t *testing.T) {
	t.Run("should return 200 when WriteHeader not called", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := nethttp_utils.NewResponseWriter(rec)

		assert.Equal(t, http.StatusOK, rw.GetStatusCode())
	})

	t.Run("should return set status after WriteHeader", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rw := nethttp_utils.NewResponseWriter(rec)
		rw.WriteHeader(http.StatusInternalServerError)

		assert.Equal(t, http.StatusInternalServerError, rw.GetStatusCode())
	})
}
