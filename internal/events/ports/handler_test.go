package ports_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"momento/internal/events/app"
	"momento/internal/events/domain"
	"momento/internal/events/mocks"
	"momento/internal/events/ports"
	"momento/pkg/listopts"
	"momento/pkg/nethttp"
	nethttp_auth "momento/pkg/nethttp/auth"
)

var mapErrorToHTTPStatus = ports.MapErrorToHTTPStatus

func TestNewEventHandler(t *testing.T) {
	t.Run("should create event handler", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		assert.NotNil(t, handler)
	})
}

func TestEventHandler_CreateEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	t.Run("should return 201 when event is created successfully", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("POST /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Event content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var got ports.EventResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.OwnerUserID)
		assert.Equal(t, "Title", got.Title)
		assert.Equal(t, "Event content", got.Content)
	})

	t.Run("should return 400 when content is empty", func(t *testing.T) {
		svc := app.NewEventService(nil)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "content cannot be empty", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /events", handler.CreateEvent)

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Event content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("POST /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Event content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})

	t.Run("should return 400 when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateEvent(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})
}

func TestEventHandler_ListEvents(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	t.Run("should return 200 with events list", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		now := time.Now().UTC()
		mockList := listopts.Paginated[domain.Event]{
			Data: []domain.Event{
				{
					ID:          primitive.NewObjectID().Hex(),
					OwnerUserID: userID,
					Content:     domain.EventContent("Event 1 content"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          primitive.NewObjectID().Hex(),
					OwnerUserID: userID,
					Content:     domain.EventContent("Event 2 content"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			Pagination: listopts.PaginationOutput{
				TotalCount: 2,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			},
		}
		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).
			Return(mockList, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListEvents(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got listopts.PaginatedResponse[ports.EventResponse]
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Len(t, got.Data, 2)
		assert.Equal(t, int64(2), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 1, got.Pagination.TotalPages)
		assert.NotEmpty(t, got.Data[0].ID)
		assert.Equal(t, "Event 1 content", got.Data[0].Content)
		assert.NotEmpty(t, got.Data[1].ID)
		assert.Equal(t, "Event 2 content", got.Data[1].Content)
	})

	t.Run("should return 200 with empty list when no events found", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mockList := listopts.Paginated[domain.Event]{
			Data: []domain.Event{},
			Pagination: listopts.PaginationOutput{
				TotalCount: 0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
			},
		}
		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).
			Return(mockList, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListEvents(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got listopts.PaginatedResponse[ports.EventResponse]
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Empty(t, got.Data)
		assert.Equal(t, int64(0), got.Pagination.TotalCount)
	})

	t.Run("should use default values when query parameters are missing", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mockList := listopts.Paginated[domain.Event]{
			Data: []domain.Event{},
			Pagination: listopts.PaginationOutput{
				TotalCount: 0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
			},
		}
		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).Return(mockList, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListEvents(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should parse query parameters correctly", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mockList := listopts.Paginated[domain.Event]{
			Data: []domain.Event{},
			Pagination: listopts.PaginationOutput{
				TotalCount: 0,
				Page:       2,
				PageSize:   10,
				TotalPages: 0,
			},
		}
		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).Return(mockList, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListEvents(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events?page=2&page_size=10&sort_by=updated_at&sort_order=asc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events", handler.ListEvents)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).
			Return(listopts.Paginated[domain.Event]{}, assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListEvents(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestEventHandler_GetUserEventByID(t *testing.T) {
	t.Run("should return 200 with event details", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		now := time.Now().UTC()
		mockEvent := domain.Event{
			ID:          primitive.NewObjectID().Hex(),
			OwnerUserID: primitive.NewObjectID().Hex(),
			Content:     domain.EventContent("Event content"),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockEvent.ID, mockEvent.OwnerUserID).Return(mockEvent, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, mockEvent.OwnerUserID)
			handler.GetUserEventByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", mockEvent.ID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got ports.EventResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, mockEvent.ID, got.ID)
		assert.Equal(t, "Event content", got.Content)
	})

	t.Run("should return 404 when event not found", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(domain.Event{}, domain.ErrEventNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.GetUserEventByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "event not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events/{id}", handler.GetUserEventByID)

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(domain.Event{}, assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.GetUserEventByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestEventHandler_UpdateEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	eventID := primitive.NewObjectID().Hex()

	t.Run("should return 200 when event is updated successfully", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		now := time.Now().UTC()
		mockEvent := domain.Event{
			ID:          eventID,
			OwnerUserID: userID,
			Content:     domain.EventContent("Event content"),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockEvent.ID, mockEvent.OwnerUserID).Return(mockEvent, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Updated title",
			"content": "Updated event content",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got ports.EventResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.OwnerUserID)
		assert.Equal(t, "Updated title", got.Title)
		assert.Equal(t, "Updated event content", got.Content)
	})

	t.Run("should return 200 when only title is updated (partial update)", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		now := time.Now().UTC()
		mockEvent := domain.Event{
			ID:          eventID,
			OwnerUserID: userID,
			Title:       domain.EventTitle("Original title"),
			Content:     domain.EventContent("Original content"),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockEvent.ID, mockEvent.OwnerUserID).Return(mockEvent, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title": "Updated title",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got ports.EventResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, "Updated title", got.Title)
		assert.Equal(t, "Original content", got.Content)
	})

	t.Run("should return 400 when content is empty", func(t *testing.T) {
		svc := app.NewEventService(nil)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "content cannot be empty", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /events/{id}", handler.UpdateEvent)

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Event content",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		now := time.Now().UTC()
		mockEvent := domain.Event{
			ID:          eventID,
			OwnerUserID: userID,
			Content:     domain.EventContent("Event content"),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(mockEvent, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateEvent(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Event content",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})

	t.Run("should return 400 when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, uri, strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "invalid request body", got.Message)
	})
}

func TestEventHandler_ArchiveEvent(t *testing.T) {
	t.Run("should return 204 no content", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().ArchiveByIDAndUserID(mock.Anything, eventID, userID).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ArchiveEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s/archive", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 404 when event not found", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().ArchiveByIDAndUserID(mock.Anything, eventID, userID).Return(domain.ErrEventNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ArchiveEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s/archive", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "event not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/archive", handler.ArchiveEvent)

		uri := fmt.Sprintf("/events/%s/archive", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().ArchiveByIDAndUserID(mock.Anything, eventID, userID).Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ArchiveEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s/archive", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestEventHandler_RestoreEvent(t *testing.T) {
	t.Run("should return 204 no content", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().RestoreByIDAndUserID(mock.Anything, eventID, userID).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.RestoreEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s/restore", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 404 when event not found", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().RestoreByIDAndUserID(mock.Anything, eventID, userID).Return(domain.ErrEventNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.RestoreEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s/restore", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "event not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/restore", handler.RestoreEvent)

		uri := fmt.Sprintf("/events/%s/restore", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().RestoreByIDAndUserID(mock.Anything, eventID, userID).Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /events/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.RestoreEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s/restore", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestEventHandler_DeleteEvent(t *testing.T) {
	t.Run("should return 204 no content", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().DeleteByIDAndUserID(mock.Anything, eventID, userID).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.DeleteEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 404 when event not found", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().DeleteByIDAndUserID(mock.Anything, eventID, userID).Return(domain.ErrEventNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.DeleteEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "event not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /events/{id}", handler.DeleteEvent)

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockEventRepository(t)
		svc := app.NewEventService(mockRepo)
		handler := ports.NewEventHandler(svc)

		eventID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().DeleteByIDAndUserID(mock.Anything, eventID, userID).Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /events/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.DeleteEvent(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/events/%s", eventID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestMapErrorToHTTPStatus(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		wantCode int
		wantMsg  string
	}{
		{"title empty", domain.ErrTitleEmpty, http.StatusBadRequest, domain.ErrTitleEmpty.Error()},
		{"title too long", domain.ErrTitleTooLong, http.StatusBadRequest, domain.ErrTitleTooLong.Error()},
		{"content empty", domain.ErrContentEmpty, http.StatusBadRequest, domain.ErrContentEmpty.Error()},
		{"content too long", domain.ErrContentTooLong, http.StatusBadRequest, domain.ErrContentTooLong.Error()},
		{"event not found", domain.ErrEventNotFound, http.StatusNotFound, domain.ErrEventNotFound.Error()},
		{"unknown error", assert.AnError, http.StatusInternalServerError, "internal server error"},
	}

	for _, tc := range testCases {
		t.Run("should handle "+tc.name, func(t *testing.T) {
			status, message := mapErrorToHTTPStatus(tc.err)
			assert.Equal(t, tc.wantCode, status)
			assert.Equal(t, tc.wantMsg, message)
		})
	}

	t.Run("should not leak internal wrapper text on unknown error", func(t *testing.T) {
		wrapped := fmt.Errorf("s.eventRepository.GetByIDAndUserID: %w", assert.AnError)
		code, msg := mapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, "internal server error", msg)
		assert.NotContains(t, msg, "eventRepository")
		assert.NotContains(t, msg, "GetByIDAndUserID")
	})

	t.Run("should use canonical domain message when sentinel is wrapped", func(t *testing.T) {
		wrapped := fmt.Errorf("app.UpdateEvent: %w", domain.ErrEventNotFound)
		code, msg := mapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusNotFound, code)
		assert.Equal(t, domain.ErrEventNotFound.Error(), msg)
		assert.NotContains(t, msg, "app.UpdateEvent")
	})
}
