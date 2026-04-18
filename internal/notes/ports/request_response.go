package ports

import (
	"time"

	"momento/internal/notes/app"
	"momento/pkg/listopts"
)

type CreateNoteRequest struct {
	Title   string `json:"title" example:"My note title"`
	Content string `json:"content" example:"My important note content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title" example:"My updated note title"`
	Content string `json:"content" example:"My updated note content"`
}

type NoteResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	UserID    string `json:"user_id" example:"507f1f77bcf86cd799439011"`
	Title     string `json:"title" example:"My note title"`
	Content   string `json:"content" example:"My important note content"`
	CreatedAt string `json:"created_at" example:"2026-02-08T10:30:00Z"`
	UpdatedAt string `json:"updated_at" example:"2026-02-08T10:30:00Z"`
}

type ListNotesResponse struct {
	Data       []NoteResponse              `json:"data"`
	Pagination listopts.PaginationResponse `json:"pagination"`
}

func NoteApplicationToResponse(note app.NoteOutput) NoteResponse {
	return NoteResponse{
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     string(note.Title),
		Content:   string(note.Content),
		CreatedAt: note.CreatedAt.Format(time.RFC3339),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
	}
}
