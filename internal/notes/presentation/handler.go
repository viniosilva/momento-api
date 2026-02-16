package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	shareddto "pinnado/internal/shared/application/dto"
	sharedresp "pinnado/internal/shared/presentation/response"
	"pinnado/pkg/nethttp"
	"pinnado/pkg/pagination"
	"pinnado/pkg/tools"
)

type noteHandler struct {
	noteService NoteService
}

func NewNoteHandler(noteService NoteService) *noteHandler {
	return &noteHandler{
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
// @Failure 400 {object} response.ErrorResponse "Invalid content"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes [post]
func (h *noteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp.JSON(w, http.StatusBadRequest, sharedresp.ErrorResponse{
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
		nethttp.JSON(w, statusCode, sharedresp.ErrorResponse{
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

// ListNotes godoc
// @Summary List user notes
// @Description Retrieves a paginated list of notes for the authenticated user with sorting options
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number (default: 1)" default(1)
// @Param page_size query int false "Page size (default: 10, max: 100)" default(10)
// @Param sort_by query string false "Sort field: created_at, updated_at (default: created_at)" default(created_at)
// @Param sort_order query string false "Sort order: asc, desc (default: desc)" default(desc)
// @Success 200 {object} ListNotesResponse
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes [get]
func (h *noteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	query := r.URL.Query()
	paginationInput := pagination.NewPaginationInput(
		tools.AtoiOrDefault(query.Get("page"), 1),
		tools.AtoiOrDefault(query.Get("page_size"), 10),
	)
	sort := shareddto.NewSortInput(
		query.Get("sort_by"),
		shareddto.OrderTypeFromString(query.Get("sort_order")),
	)

	input := application.ListNotesInput{
		UserID:     userID,
		Pagination: paginationInput,
		Sort:       sort,
	}

	output, err := h.noteService.ListNotes(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	data := make([]NoteResponse, 0, len(output.Data))
	for _, note := range output.Data {
		data = append(data, NoteApplicationToResponse(note))
	}

	response := ListNotesResponse{
		Data:       data,
		Pagination: pagination.PaginationApplicationToResponse(output.Pagination),
	}

	nethttp.JSON(w, http.StatusOK, response)
}

func MapErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, domain.ErrContentEmpty) ||
		errors.Is(err, domain.ErrContentTooLong) ||
		errors.Is(err, domain.ErrInvalidNoteContent) {
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusInternalServerError, "internal server error"
}
