package presentation

import (
	"encoding/json"
	"errors"
	"momento/internal/auth/application"
	"momento/internal/auth/domain"
	nethttp_utils "momento/pkg/nethttp/utils"
	"net/http"
	"strings"
)

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{
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
// @Router /api/auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, ErrorResponse{
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
		nethttp_utils.JSON(w, statusCode, ErrorResponse{
			Message: message,
		})
		return
	}

	response := RegisterResponse{
		ID:    output.ID,
		Email: string(output.Email),
	}

	nethttp_utils.JSON(w, http.StatusCreated, response)
}

// Login godoc
// @Summary Login with email and password
// @Description Authenticates a user with email and password credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/auth/login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := application.LoginInput{
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
	}

	output, err := h.authService.Login(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, ErrorResponse{
			Message: message,
		})
		return
	}

	response := LoginResponse{
		Token: output.Token,
	}

	nethttp_utils.JSON(w, http.StatusOK, response)
}

// MapErrorToHTTPStatus maps domain/application errors to appropriate HTTP status codes
// Exported for testing purposes
func MapErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, domain.ErrUserAlreadyExists) {
		return http.StatusConflict, err.Error()
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

	if errors.Is(err, domain.ErrInvalidCredentials) {
		return http.StatusUnauthorized, err.Error()
	}

	return http.StatusInternalServerError, "internal server error"
}
