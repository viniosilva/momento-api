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

	"momento/internal/notes/app"
	"momento/internal/notes/domain"
	"momento/internal/notes/mocks"
	"momento/internal/notes/ports"
	"momento/pkg/listopts"
	"momento/pkg/nethttp"
	nethttp_auth "momento/pkg/nethttp/auth"
)

var mapErrorToHTTPStatus = ports.MapErrorToHTTPStatus

func TestNewNoteHandler(t *testing.T) {
	t.Run("should create note handler", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		assert.NotNil(t, handler)
	})
}

func TestNoteHandler_CreateNote(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	t.Run("should return 201 when note is created successfully", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Note content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var got ports.NoteResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, "Title", got.Title)
		assert.Equal(t, "Note content", got.Content)
	})

	t.Run("should return 400 when content is empty", func(t *testing.T) {
		svc := app.NewNoteService(nil)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/notes", bytes.NewReader(body))
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /notes", handler.CreateNote)

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Note content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/notes", bytes.NewReader(body))
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Note content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/notes", bytes.NewReader(body))
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.CreateNote(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/notes", strings.NewReader("invalid json"))
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

func TestNoteHandler_ListNotes(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	t.Run("should return 200 with notes list", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		now := time.Now().UTC()
		mockList := listopts.Paginated[domain.Note]{
			Data: []domain.Note{
				{
					ID:        primitive.NewObjectID().Hex(),
					UserID:    userID,
					Content:   domain.NoteContent("Note 1 content"),
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        primitive.NewObjectID().Hex(),
					UserID:    userID,
					Content:   domain.NoteContent("Note 2 content"),
					CreatedAt: now,
					UpdatedAt: now,
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
		mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListNotes(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got listopts.PaginatedResponse[ports.NoteResponse]
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Len(t, got.Data, 2)
		assert.Equal(t, int64(2), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 1, got.Pagination.TotalPages)
		assert.NotEmpty(t, got.Data[0].ID)
		assert.Equal(t, "Note 1 content", got.Data[0].Content)
		assert.NotEmpty(t, got.Data[1].ID)
		assert.Equal(t, "Note 2 content", got.Data[1].Content)
	})

	t.Run("should return 200 with empty list when no notes found", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mockList := listopts.Paginated[domain.Note]{
			Data: []domain.Note{},
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
		mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListNotes(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got listopts.PaginatedResponse[ports.NoteResponse]
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Empty(t, got.Data)
		assert.Equal(t, int64(0), got.Pagination.TotalCount)
	})

	t.Run("should use default values when query parameters are missing", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mockList := listopts.Paginated[domain.Note]{
			Data: []domain.Note{},
			Pagination: listopts.PaginationOutput{
				TotalCount: 0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
			},
		}
		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).Return(mockList, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListNotes(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should parse query parameters correctly", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mockList := listopts.Paginated[domain.Note]{
			Data: []domain.Note{},
			Pagination: listopts.PaginationOutput{
				TotalCount: 0,
				Page:       2,
				PageSize:   10,
				TotalPages: 0,
			},
		}
		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).Return(mockList, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListNotes(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes?page=2&page_size=10&sort_by=updated_at&sort_order=asc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes", handler.ListNotes)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).
			Return(listopts.Paginated[domain.Note]{}, assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ListNotes(w, r.WithContext(ctx))
		})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/notes", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestNoteHandler_GetUserNoteByID(t *testing.T) {
	t.Run("should return 200 with note details", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		now := time.Now().UTC()
		mockNote := domain.Note{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    primitive.NewObjectID().Hex(),
			Content:   domain.NoteContent("Note content"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockNote.ID, mockNote.UserID).Return(mockNote, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, mockNote.UserID)
			handler.GetUserNoteByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", mockNote.ID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got ports.NoteResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, mockNote.ID, got.ID)
		assert.Equal(t, "Note content", got.Content)
	})

	t.Run("should return 404 when note not found", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(domain.Note{}, domain.ErrNoteNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.GetUserNoteByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "note not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", handler.GetUserNoteByID)

		uri := fmt.Sprintf("/notes/%s", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(domain.Note{}, assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.GetUserNoteByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID)
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

func TestNoteHandler_UpdateNote(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	noteID := primitive.NewObjectID().Hex()

	t.Run("should return 200 when note is updated successfully", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		now := time.Now().UTC()
		mockNote := domain.Note{
			ID:        noteID,
			UserID:    userID,
			Content:   domain.NoteContent("Note content"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockNote.ID, mockNote.UserID).Return(mockNote, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Updated title",
			"content": "Updated note content",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/notes/%s", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got ports.NoteResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, "Updated title", got.Title)
		assert.Equal(t, "Updated note content", got.Content)
	})

	t.Run("should return 200 when only title is updated (partial update)", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		now := time.Now().UTC()
		mockNote := domain.Note{
			ID:        noteID,
			UserID:    userID,
			Title:     domain.NoteTitle("Original title"),
			Content:   domain.NoteContent("Original content"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockNote.ID, mockNote.UserID).Return(mockNote, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title": "Updated title",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/notes/%s", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got ports.NoteResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, "Updated title", got.Title)
		assert.Equal(t, "Original content", got.Content)
	})

	t.Run("should return 400 when content is empty", func(t *testing.T) {
		svc := app.NewNoteService(nil)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/notes/%s", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /notes/{id}", handler.UpdateNote)

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Note content",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/notes/%s", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		now := time.Now().UTC()
		mockNote := domain.Note{
			ID:        noteID,
			UserID:    userID,
			Content:   domain.NoteContent("Note content"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(mockNote, nil).Once()
		mockRepo.EXPECT().Update(mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateNote(w, r.WithContext(ctx))
		})

		reqBody := map[string]any{
			"title":   "Title",
			"content": "Note content",
		}
		body, _ := json.Marshal(reqBody)
		uri := fmt.Sprintf("/notes/%s", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.UpdateNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID)
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

func TestNoteHandler_ArchiveNote(t *testing.T) {
	t.Run("should return 204 no content", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().ArchiveByIDAndUserID(mock.Anything, noteID, userID).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ArchiveNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s/archive", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 404 when note not found", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().ArchiveByIDAndUserID(mock.Anything, noteID, userID).Return(domain.ErrNoteNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ArchiveNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s/archive", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "note not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/archive", handler.ArchiveNote)

		uri := fmt.Sprintf("/notes/%s/archive", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().ArchiveByIDAndUserID(mock.Anything, noteID, userID).Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.ArchiveNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s/archive", noteID)
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

func TestNoteHandler_RestoreNote(t *testing.T) {
	t.Run("should return 204 no content", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().RestoreByIDAndUserID(mock.Anything, noteID, userID).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.RestoreNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s/restore", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 404 when note not found", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().RestoreByIDAndUserID(mock.Anything, noteID, userID).Return(domain.ErrNoteNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.RestoreNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s/restore", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "note not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/restore", handler.RestoreNote)

		uri := fmt.Sprintf("/notes/%s/restore", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().RestoreByIDAndUserID(mock.Anything, noteID, userID).Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("PATCH /notes/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.RestoreNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s/restore", noteID)
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

func TestNoteHandler_DeleteNote(t *testing.T) {
	t.Run("should return 204 no content", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().DeleteByIDAndUserID(mock.Anything, noteID, userID).Return(nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.DeleteNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should return 404 when note not found", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().DeleteByIDAndUserID(mock.Anything, noteID, userID).Return(domain.ErrNoteNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.DeleteNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got nethttp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "note not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /notes/{id}", handler.DeleteNote)

		uri := fmt.Sprintf("/notes/%s", noteID)
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
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := app.NewNoteService(mockRepo)
		handler := ports.NewNoteHandler(svc)

		noteID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		mockRepo.EXPECT().DeleteByIDAndUserID(mock.Anything, noteID, userID).Return(assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID)
			handler.DeleteNote(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID)
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
		{"note not found", domain.ErrNoteNotFound, http.StatusNotFound, domain.ErrNoteNotFound.Error()},
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
		wrapped := fmt.Errorf("s.noteRepository.GetByIDAndUserID: %w", assert.AnError)
		code, msg := mapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, "internal server error", msg)
		assert.NotContains(t, msg, "noteRepository")
		assert.NotContains(t, msg, "GetByIDAndUserID")
	})

	t.Run("should use canonical domain message when sentinel is wrapped", func(t *testing.T) {
		wrapped := fmt.Errorf("app.UpdateNote: %w", domain.ErrNoteNotFound)
		code, msg := mapErrorToHTTPStatus(wrapped)

		assert.Equal(t, http.StatusNotFound, code)
		assert.Equal(t, domain.ErrNoteNotFound.Error(), msg)
		assert.NotContains(t, msg, "app.UpdateNote")
	})
}
