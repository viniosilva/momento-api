package nethttp

// ErrorResponse is the standard error body returned on all HTTP error responses.
type ErrorResponse struct {
	Message string `json:"message" example:"internal server error"`
}
