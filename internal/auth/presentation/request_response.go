package presentation

type RegisterRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"ValidPass123!"`
}

type RegisterResponse struct {
	ID    string `json:"id" example:"507f1f77bcf86cd799439011"`
	Email string `json:"email" example:"user@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"ValidPass123!"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Message string `json:"message" example:"user already exists"`
}
