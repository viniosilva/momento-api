package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"momento/internal/events/app"
	"momento/internal/events/domain"
	"momento/pkg/listopts"
	"momento/pkg/nethttp"
	nethttp_auth "momento/pkg/nethttp/auth"
	nethttp_utils "momento/pkg/nethttp/utils"
	"momento/pkg/tools"
)

type eventHandler struct {
	eventService EventService
}

func NewEventHandler(eventService EventService) *eventHandler {
	return &eventHandler{
		eventService: eventService,
	}
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Creates a new event associated with the authenticated user
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateEventRequest true "Event content"
// @Success 201 {object} EventResponse
// @Failure 400 {object} nethttp.ErrorResponse "Invalid content"
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events [post]
func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := app.EventInput{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	output, err := h.eventService.CreateEvent(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := EventResponse{
		ID:          output.ID,
		OwnerUserID: output.OwnerUserID,
		Title:       string(output.Title),
		Content:     string(output.Content),
		CreatedAt:   output.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   output.UpdatedAt.Format(time.RFC3339),
	}

	nethttp_utils.JSON(w, http.StatusCreated, res)
}

// ListEvents godoc
// @Summary List user events
// @Description Retrieves a paginated list of events for the authenticated user with sorting options
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number (default: 1)" default(1)
// @Param page_size query int false "Page size (default: 10, max: 100)" default(10)
// @Param sort_by query string false "Sort field: created_at, updated_at (default: created_at)" default(created_at)
// @Param sort_order query string false "Sort order: asc, desc (default: desc)" default(desc)
// @Success 200 {object} ListEventsResponse
// @Failure 400 {object} nethttp.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events [get]
func (h *eventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
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

	input := app.ListEventsInput{
		UserID:     userID,
		Pagination: paginationInput,
		Sort:       sort,
	}

	output, err := h.eventService.ListEvents(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	data := make([]EventResponse, 0, len(output.Data))
	for _, event := range output.Data {
		data = append(data, EventApplicationToResponse(event))
	}

	res := ListEventsResponse{
		Data:       data,
		Pagination: listopts.PaginationApplicationToResponse(output.Pagination),
	}

	nethttp_utils.JSON(w, http.StatusOK, res)
}

// GetUserEventByID godoc
// @Summary Retrieve a event by ID
// @Description Retrieves a specific event for the authenticated user by event ID
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Success 200 {object} EventResponse
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} nethttp.ErrorResponse "Event not found"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events/{id} [get]
func (h *eventHandler) GetUserEventByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := app.GetUserEventByIDInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	output, err := h.eventService.GetUserEventByID(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := EventApplicationToResponse(output)
	nethttp_utils.JSON(w, http.StatusOK, res)
}

// UpdateEvent godoc
// @Summary Update a event
// @Description Updates the content of a event for the authenticated user
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Param request body UpdateEventRequest true "Updated event content"
// @Success 200 {object} EventResponse
// @Failure 400 {object} nethttp.ErrorResponse "Invalid content"
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} nethttp.ErrorResponse "Event not found"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events/{id} [patch]
func (h *eventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, nethttp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := app.UpdateEventInput{
		UserID:  userID,
		ID:      r.PathValue("id"),
		Title:   req.Title,
		Content: req.Content,
	}

	output, err := h.eventService.UpdateEvent(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := EventApplicationToResponse(output)
	nethttp_utils.JSON(w, http.StatusOK, res)
}

// ArchiveEvent godoc
// @Summary Archive a event
// @Description Archives a event belonging to the authenticated user
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Success 204
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} nethttp.ErrorResponse "Event not found"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events/{id}/archive [patch]
func (h *eventHandler) ArchiveEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := app.ArchiveEventInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	err := h.eventService.ArchiveEvent(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

// RestoreEvent godoc
// @Summary Restore a event
// @Description Restores a event belonging to the authenticated user
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Success 204
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} nethttp.ErrorResponse "Event not found"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events/{id}/restore [patch]
func (h *eventHandler) RestoreEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := app.RestoreEventInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	err := h.eventService.RestoreEvent(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	nethttp_utils.StatusCode(w, http.StatusNoContent)
}

// DeleteEvent godoc
// @Summary Delete a event
// @Description Deletes a event belonging to the authenticated user
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Success 204
// @Failure 401 {object} nethttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} nethttp.ErrorResponse "Event not found"
// @Failure 500 {object} nethttp.ErrorResponse "Internal server error"
// @Router /api/events/{id} [delete]
func (h *eventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(nethttp_auth.ContextKeyUserID).(string)
	if !ok || userID == "" {
		nethttp_utils.JSON(w, http.StatusUnauthorized, nethttp.ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	input := app.DeleteEventInput{
		UserID: userID,
		ID:     r.PathValue("id"),
	}

	if err := h.eventService.DeleteEvent(r.Context(), input); err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, nethttp.ErrorResponse{
			Message: message,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func MapErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrTitleEmpty):
		return http.StatusBadRequest, domain.ErrTitleEmpty.Error()
	case errors.Is(err, domain.ErrTitleTooLong):
		return http.StatusBadRequest, domain.ErrTitleTooLong.Error()
	case errors.Is(err, domain.ErrContentEmpty):
		return http.StatusBadRequest, domain.ErrContentEmpty.Error()
	case errors.Is(err, domain.ErrContentTooLong):
		return http.StatusBadRequest, domain.ErrContentTooLong.Error()
	case errors.Is(err, domain.ErrEventNotFound):
		return http.StatusNotFound, domain.ErrEventNotFound.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
