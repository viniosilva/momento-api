package nethttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type callback func(w http.ResponseWriter, r *http.Request)

func Request[T any](ctx context.Context, method string, target string, body T, cb callback) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequestWithContext(ctx, method, target, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	cb(rec, req)

	return rec.Result(), nil
}

// RequestWithResponse makes a request to the given target with the given method and body and returns the response and the response body.
// T is the type of the response body.
// R is the type of the request body.
func RequestWithResponse[T any, R any](ctx context.Context, method string, target string, body T, cb callback) (*http.Response, *R, error) {
	resp, err := Request(ctx, method, target, body, cb)
	if err != nil {
		return nil, nil, err
	}

	var respR R
	err = json.NewDecoder(resp.Body).Decode(&respR)
	if err != nil {
		return nil, nil, err
	}

	return resp, &respR, nil
}
