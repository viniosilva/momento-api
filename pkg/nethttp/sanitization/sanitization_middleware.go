package nethttp_sanitization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	nethttp_utils "pinnado/pkg/nethttp/utils"

	"github.com/microcosm-cc/bluemonday"
)

const (
	maxBodySize      = 1 << 20 // 1MB
	maxStringSize    = 5120    // 5KB
	maxKeySize       = 100     // Key max length
	maxNestedDepth   = 10      // Max recursion depth
	allowedMimeTypes = "application/json"
)

func SanitizationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 || r.Method == http.MethodGet || r.Method == http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, allowedMimeTypes) {
				next.ServeHTTP(w, r)
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

			body, err := io.ReadAll(r.Body)
			if err != nil {
				nethttp_utils.JSON(w, http.StatusBadRequest, map[string]string{
					"message": "request body too large (max 1MB)",
				})
				return
			}

			var data any
			if err := json.Unmarshal(body, &data); err != nil {
				nethttp_utils.JSON(w, http.StatusBadRequest, map[string]string{
					"message": "invalid JSON format",
				})
				return
			}
			sanitized, err := sanitizeValue(data, 0)
			if err != nil {
				nethttp_utils.JSON(w, http.StatusBadRequest, map[string]string{
					"message": fmt.Sprintf("invalid input: %s", err.Error()),
				})
				return
			}

			sanitizedBody, err := json.Marshal(sanitized)
			if err != nil {
				nethttp_utils.JSON(w, http.StatusInternalServerError, map[string]string{
					"message": "internal server error",
				})
				return
			}
			r.Body = io.NopCloser(strings.NewReader(string(sanitizedBody)))

			next.ServeHTTP(w, r)
		})
	}
}
func sanitizeValue(data any, depth int) (any, error) {
	if depth > maxNestedDepth {
		return nil, fmt.Errorf("nesting too deep (max %d levels)", maxNestedDepth)
	}

	switch v := data.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			if utf8.RuneCountInString(key) > maxKeySize {
				return nil, fmt.Errorf("key too long (max %d chars): %s", maxKeySize, key)
			}

			sanitized, err := sanitizeValue(value, depth+1)
			if err != nil {
				return nil, err
			}
			result[key] = sanitized
		}
		return result, nil

	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			sanitized, err := sanitizeValue(item, depth+1)
			if err != nil {
				return nil, err
			}
			result[i] = sanitized
		}
		return result, nil

	case string:
		return sanitizeString(v)

	case float64, bool, nil:
		return v, nil

	default:
		return v, nil
	}
}

func sanitizeString(s string) (string, error) {
	runeCount := utf8.RuneCountInString(s)
	if runeCount > maxStringSize {
		return "", fmt.Errorf("string too long (max %d chars)", maxStringSize)
	}
	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(s)
	sanitized = removeControlCharacters(sanitized)

	return sanitized, nil
}

func removeControlCharacters(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' && r != '\r' {
			return -1 // Remove
		}
		if r == 127 {
			return -1
		}

		return r
	}, s)
}
