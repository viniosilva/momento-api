package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"momento/internal/auth/app"
	"momento/internal/auth/domain"
	"momento/pkg/nethttp"
	nethttp_utils "momento/pkg/nethttp/utils"
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
// @Failure 400 {object} nethttp.ErrorResponse "Invalid request data"
// @Failure 409 {object} nethttp.ErrorResponse "User already exists"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := app.UserInput{
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
	}

	output, err := h.authService.Register(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
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
// @Description Authenticates a user and returns a JWT token and a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} nethttp.ErrorResponse "Invalid request data"
// @Failure 401 {object} nethttp.ErrorResponse "Invalid credentials"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/auth/login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := app.LoginInput{
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
	}

	output, err := h.authService.Login(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	response := LoginResponse{
		Token:        output.Token,
		RefreshToken: output.RefreshToken,
	}

	nethttp_utils.JSON(w, http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Exchanges a refresh token for a new JWT token and a new refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh request"
// @Success 200 {object} RefreshResponse "Tokens refreshed successfully"
// @Failure 400 {object} nethttp.ErrorResponse "Invalid request body"
// @Failure 401 {object} nethttp.ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/auth/refresh [post]
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "refresh_token is required",
		})
		return
	}

	output, err := h.authService.RefreshToken(r.Context(), app.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	response := RefreshResponse{
		Token:        output.Token,
		RefreshToken: output.RefreshToken,
	}

	nethttp_utils.JSON(w, http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidates the refresh token to log out the user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Logout request"
// @Success 204 "Logout successful"
// @Failure 400 {object} nethttp.ErrorResponse "Invalid request body"
// @Router /api/auth/logout [post]
func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "refresh_token is required",
		})
		return
	}

	err := h.authService.Logout(r.Context(), app.LogoutInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Sends a password reset email if the account exists
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Password reset request"
// @Success 200 {object} map[string]string "Reset email sent"
// @Failure 400 {object} nethttp.ErrorResponse "Invalid request data"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/auth/forgot-password [post]
func (h *authHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := app.ForgotPasswordInput{
		Email: strings.TrimSpace(req.Email),
	}

	if err := h.authService.ForgotPassword(r.Context(), input); err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Resets the user's password using a valid reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset password request"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} nethttp.ErrorResponse "Invalid token or password"
// @Failure 410 {object} nethttp.ErrorResponse "Token expired"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/auth/reset-password [post]
func (h *authHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := app.ResetPasswordInput{
		Token:    req.Token,
		Password: req.Password,
	}

	if err := h.authService.ResetPassword(r.Context(), input); err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

// ValidateResetToken godoc
// @Summary Validate reset token
// @Description Checks if a reset token is valid and not expired
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Reset token"
// @Success 200 {object} map[string]bool "Token is valid"
// @Failure 400 {object} nethttp.ErrorResponse "Invalid token"
// @Failure 410 {object} nethttp.ErrorResponse "Token expired"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/auth/reset-password/validate [get]
func (h *authHandler) ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if strings.TrimSpace(token) == "" {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "token is required",
		})
		return
	}

	input := app.ValidateResetTokenInput{
		Token: token,
	}

	_, err := h.authService.ValidateResetToken(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.JSON(w, http.StatusOK, ValidateResetTokenResponse{
		Valid: true,
	})
}

func MapErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, domain.ErrUserAlreadyExists.Error()
	case errors.Is(err, domain.ErrEmailIsEmpty):
		return http.StatusBadRequest, domain.ErrEmailIsEmpty.Error()
	case errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest, domain.ErrInvalidEmail.Error()
	case errors.Is(err, domain.ErrPasswordTooShort):
		return http.StatusBadRequest, domain.ErrPasswordTooShort.Error()
	case errors.Is(err, domain.ErrPasswordTooLong):
		return http.StatusBadRequest, domain.ErrPasswordTooLong.Error()
	case errors.Is(err, domain.ErrPasswordMissingUpper):
		return http.StatusBadRequest, domain.ErrPasswordMissingUpper.Error()
	case errors.Is(err, domain.ErrPasswordMissingLower):
		return http.StatusBadRequest, domain.ErrPasswordMissingLower.Error()
	case errors.Is(err, domain.ErrPasswordMissingNumber):
		return http.StatusBadRequest, domain.ErrPasswordMissingNumber.Error()
	case errors.Is(err, domain.ErrPasswordMissingSymbol):
		return http.StatusBadRequest, domain.ErrPasswordMissingSymbol.Error()
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, domain.ErrInvalidCredentials.Error()
	case errors.Is(err, domain.ErrRefreshTokenInvalid):
		return http.StatusUnauthorized, domain.ErrRefreshTokenInvalid.Error()
	case errors.Is(err, domain.ErrRefreshTokenNotFound):
		return http.StatusUnauthorized, domain.ErrRefreshTokenInvalid.Error()
	case errors.Is(err, domain.ErrRefreshTokenExpired):
		return http.StatusUnauthorized, domain.ErrRefreshTokenInvalid.Error()
	case errors.Is(err, domain.ErrInvalidResetToken):
		return http.StatusBadRequest, domain.ErrInvalidResetToken.Error()
	case errors.Is(err, domain.ErrExpiredResetToken):
		return http.StatusGone, domain.ErrExpiredResetToken.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
