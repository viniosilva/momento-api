package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pinnado/internal/auth/application"
	"pinnado/internal/auth/domain"
	"pinnado/pkg/nethttp"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} RegisterResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := application.UserInput{
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
	}

	output, err := h.authService.Register(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp.JSON(w, statusCode, ErrorResponse{
			Message: message,
		})
		return
	}

	response := RegisterResponse{
		ID:    output.ID,
		Email: string(output.Email),
	}

	nethttp.JSON(w, http.StatusCreated, response)
}

// MapErrorToHTTPStatus maps domain/application errors to appropriate HTTP status codes
// Exported for testing purposes
func MapErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, domain.ErrUserAlreadyExists) {
		return http.StatusConflict, "user already exists"
	}

	if errors.Is(err, domain.ErrEmailIsEmpty) ||
		errors.Is(err, domain.ErrInvalidEmail) ||
		errors.Is(err, domain.ErrPasswordTooShort) ||
		errors.Is(err, domain.ErrPasswordTooLong) ||
		errors.Is(err, domain.ErrPasswordMissingUpper) ||
		errors.Is(err, domain.ErrPasswordMissingLower) ||
		errors.Is(err, domain.ErrPasswordMissingNumber) ||
		errors.Is(err, domain.ErrPasswordMissingSymbol) {
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusInternalServerError, "internal server error"
}
