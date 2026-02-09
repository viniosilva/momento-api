package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	"pinnado/pkg/nethttp"
)

type NoteHandler struct {
	noteService NoteService
}

func NewNoteHandler(noteService NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

// CreateNote godoc
// @Summary Create a new note
// @Description Creates a new note associated with the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateNoteRequest true "Note content"
// @Success 201 {object} NoteResponse
// @Failure 400 {object} ErrorResponse "Invalid content"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp.JSON(w, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := application.NoteInput{
		UserID:  userID,
		Content: req.Content,
	}

	output, err := h.noteService.CreateNote(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp.JSON(w, statusCode, ErrorResponse{
			Message: message,
		})
		return
	}

	response := NoteResponse{
		ID:        output.ID,
		UserID:    output.UserID,
		Content:   string(output.Content),
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
		UpdatedAt: output.UpdatedAt.Format(time.RFC3339),
	}

	nethttp.JSON(w, http.StatusCreated, response)
}

func MapErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, domain.ErrContentEmpty) ||
		errors.Is(err, domain.ErrContentTooLong) ||
		errors.Is(err, domain.ErrInvalidNoteContent) {
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusInternalServerError, "internal server error"
}
