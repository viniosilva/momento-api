package ports

import (
	"time"

	"momento/internal/events/app"
	"momento/pkg/listopts"
)

type CreateEventRequest struct {
	Title   string `json:"title" example:"My event title"`
	Content string `json:"content" example:"My important event content"`
}

type UpdateEventRequest struct {
	Title   *string `json:"title" example:"My updated event title"`
	Content *string `json:"content" example:"My updated event content"`
}

type EventResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	UserID    string `json:"user_id" example:"507f1f77bcf86cd799439011"`
	Title     string `json:"title" example:"My event title"`
	Content   string `json:"content" example:"My important event content"`
	CreatedAt string `json:"created_at" example:"2026-02-08T10:30:00Z"`
	UpdatedAt string `json:"updated_at" example:"2026-02-08T10:30:00Z"`
}

type ListEventsResponse struct {
	Data       []EventResponse            `json:"data"`
	Pagination listopts.PaginationResponse `json:"pagination"`
}

func EventApplicationToResponse(event app.EventOutput) EventResponse {
	return EventResponse{
		ID:        event.ID,
		UserID:    event.UserID,
		Title:     string(event.Title),
		Content:   string(event.Content),
		CreatedAt: event.CreatedAt.Format(time.RFC3339),
		UpdatedAt: event.UpdatedAt.Format(time.RFC3339),
	}
}