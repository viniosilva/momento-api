package presentation

type CreateNoteRequest struct {
	Content string `json:"content" example:"My important note content"`
}

type NoteResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	UserID    string `json:"user_id" example:"507f1f77bcf86cd799439011"`
	Content   string `json:"content" example:"My important note content"`
	CreatedAt string `json:"created_at" example:"2026-02-08T10:30:00Z"`
	UpdatedAt string `json:"updated_at" example:"2026-02-08T10:30:00Z"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"invalid note content"`
}
