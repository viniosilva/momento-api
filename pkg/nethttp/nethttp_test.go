package nethttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"pinnado/pkg/nethttp"
)

func TestJSON(t *testing.T) {
	t.Run("should return 200 status with json response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{
			"result": "ok",
		}

		nethttp.JSON(w, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)

		assert.NoError(t, err)
		assert.Equal(t, "ok", got["result"])
	})
}
