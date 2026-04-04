package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	sharedresp "pinnado/internal/shared/presentation/response"
	"pinnado/pkg/listopts"
	nethttp_auth "pinnado/pkg/nethttp/auth"
	nethttp_utils "pinnado/pkg/nethttp/utils"
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
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, sharedresp.ErrorResponse{
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
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := NoteResponse{
		ID:        output.ID,
		UserID:    output.UserID,
		Content:   string(output.Content),
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
		UpdatedAt: output.UpdatedAt.Format(time.RFC3339),
	}

	nethttp_utils.JSON(w, http.StatusCreated, res)
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
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	query := r.URL.Query()
	paginationInput := listopts.NewPaginationInput(
		tools.AtoiOrDefault(query.Get("page"), 1),
		tools.AtoiOrDefault(query.Get("page_size"), 10),
	)
	sort := listopts.NewSortInput(
		query.Get("sort_by"),
		listopts.OrderTypeFromString(query.Get("sort_order")),
	)

	input := application.ListNotesInput{
		UserID:     userID,
		Pagination: paginationInput,
		Sort:       sort,
	}

	output, err := h.noteService.ListNotes(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	data := make([]NoteResponse, 0, len(output.Data))
	for _, note := range output.Data {
		data = append(data, NoteApplicationToResponse(note))
	}

	res := ListNotesResponse{
		Data:       data,
		Pagination: listopts.PaginationApplicationToResponse(output.Pagination),
	}

	nethttp_utils.JSON(w, http.StatusOK, res)
}

// GetUserNoteByID godoc
// @Summary Retrieve a note by ID
// @Description Retrieves a specific note for the authenticated user by note ID
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Note ID"
// @Success 200 {object} NoteResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Note not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes/{id} [get]
func (h *noteHandler) GetUserNoteByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := application.GetUserNoteByIDInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	output, err := h.noteService.GetUserNoteByID(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := NoteApplicationToResponse(output)
	nethttp_utils.JSON(w, http.StatusOK, res)
}

// UpdateNote godoc
// @Summary Update a note
// @Description Updates the content of a note for the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Note ID"
// @Param request body UpdateNoteRequest true "Updated note content"
// @Success 200 {object} NoteResponse
// @Failure 400 {object} response.ErrorResponse "Invalid content"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Note not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes/{id} [patch]
func (h *noteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, sharedresp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := application.UpdateNoteInput{
		UserID:  userID,
		ID:      r.PathValue("id"),
		Content: req.Content,
	}

	output, err := h.noteService.UpdateNote(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := NoteApplicationToResponse(output)
	nethttp_utils.JSON(w, http.StatusOK, res)
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Deletes a note belonging to the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Note ID"
// @Success 204
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Note not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes/{id} [delete]
func (h *noteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := application.DeleteNoteInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	if err := h.noteService.DeleteNote(r.Context(), input); err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ArchiveNote godoc
// @Summary Archive a note
// @Description Archives a note belonging to the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Note ID"
// @Success 200 {object} NoteResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Note not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes/{id}/archive [patch]
func (h *noteHandler) ArchiveNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := application.ArchiveNoteInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	err := h.noteService.ArchiveNote(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

// RestoreNote godoc
// @Summary Restore a note
// @Description Restores a note belonging to the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Note ID"
// @Success 200 {object} NoteResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Note not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/notes/{id}/restore [patch]
func (h *noteHandler) RestoreNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, sharedresp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := application.RestoreNoteInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	err := h.noteService.RestoreNote(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

func MapErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, domain.ErrContentEmpty) ||
		errors.Is(err, domain.ErrContentTooLong) ||
		errors.Is(err, domain.ErrInvalidNoteContent) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, domain.ErrNoteNotFound) {
		return http.StatusNotFound, err.Error()
	}

	return http.StatusInternalServerError, "internal server error"
}
