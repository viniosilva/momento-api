package presentation_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	"pinnado/internal/notes/mocks"
	"pinnado/internal/notes/presentation"
	sharedresp "pinnado/internal/shared/presentation/response"
	"pinnado/pkg/listopts"
	"pinnado/pkg/nethttp"
	nethttp_auth "pinnado/pkg/nethttp/auth"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mapErrorToHTTPStatus = presentation.MapErrorToHTTPStatus

func TestNewNoteHandler(t *testing.T) {
	t.Run("should create note handler", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		assert.NotNil(t, handler)
	})
}

func TestNoteHandler_CreateNote(t *testing.T) {
	validUserID := "507f1f77bcf86cd799439011"

	t.Run("should return 201 when note is created successfully", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
			Return(nil).
			Once()

		reqBody := map[string]any{
			"content": "Valid note content",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.NoteResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, validUserID, got.UserID)
		assert.Equal(t, "Valid note content", got.Content)
	})

	t.Run("should return 400 when content is empty", func(t *testing.T) {
		svc := application.NewNoteService(nil)
		handler := presentation.NewNoteHandler(svc)

		reqBody := map[string]any{
			"content": "",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "content cannot be empty", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		reqBody := map[string]any{
			"content": "Valid note content",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, handler.CreateNote)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		reqBody := map[string]any{
			"content": "Valid note content",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})

	t.Run("should return 400 when request body is invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[string, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", "invalid json", func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid request body", got.Message)
	})
}

func TestMapErrorToHTTPStatus(t *testing.T) {
	t.Run("should return 400 for content empty error", func(t *testing.T) {
		status, message := mapErrorToHTTPStatus(domain.ErrContentEmpty)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, domain.ErrContentEmpty.Error(), message)
	})

	t.Run("should return 400 for content too long error", func(t *testing.T) {
		status, message := mapErrorToHTTPStatus(domain.ErrContentTooLong)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, domain.ErrContentTooLong.Error(), message)
	})

	t.Run("should return 400 for invalid note content error", func(t *testing.T) {
		status, message := mapErrorToHTTPStatus(domain.ErrInvalidNoteContent)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, domain.ErrInvalidNoteContent.Error(), message)
	})

	t.Run("should return 500 for unknown error", func(t *testing.T) {
		status, message := mapErrorToHTTPStatus(assert.AnError)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "internal server error", message)
	})
}

func TestNoteHandler_ListNotes(t *testing.T) {
	validUserID := "507f1f77bcf86cd799439011"

	t.Run("should return 200 with notes list", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		userID := primitive.NewObjectID()
		now := time.Now().UTC()
		mockList := listopts.Paginated[domain.Note]{
			Data: []domain.Note{
				{
					ID:        primitive.NewObjectID(),
					UserID:    userID,
					Content:   domain.NoteContent("Note 1 content"),
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        primitive.NewObjectID(),
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
			Return(mockList, nil).
			Once()

		resp, got, err := nethttp.RequestWithResponse[any, listopts.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes?page=1&page_size=10", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

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
			Return(mockList, nil).
			Once()

		resp, got, err := nethttp.RequestWithResponse[any, listopts.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Empty(t, got.Data)
		assert.Equal(t, int64(0), got.Pagination.TotalCount)
	})

	t.Run("should use default values when query parameters are missing", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

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

		resp, _, err := nethttp.RequestWithResponse[any, listopts.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should parse query parameters correctly", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

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

		resp, _, err := nethttp.RequestWithResponse[any, listopts.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes?page=2&page_size=10&sort_by=updated_at&sort_order=asc", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		resp, got, err := nethttp.RequestWithResponse[any, sharedresp.ErrorResponse](
			t.Context(), http.MethodGet, "/notes", nil, handler.ListNotes)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		mockRepo.EXPECT().ListByUserID(mock.Anything, mock.Anything, mock.Anything).
			Return(listopts.Paginated[domain.Note]{}, assert.AnError).
			Once()

		resp, got, err := nethttp.RequestWithResponse[any, sharedresp.ErrorResponse](
			t.Context(), http.MethodGet, "/notes", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}

func TestNoteHandler_GetUserNoteByID(t *testing.T) {
	t.Run("should return 200 with note details", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		now := time.Now().UTC()
		mockNote := domain.Note{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Content:   domain.NoteContent("Note content"),
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, mockNote.ID, mockNote.UserID).Return(mockNote, nil).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, mockNote.UserID.Hex())
			handler.GetUserNoteByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", mockNote.ID.Hex())
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got presentation.NoteResponse
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)

		assert.Equal(t, mockNote.ID.Hex(), got.ID)
		assert.Equal(t, "Note content", got.Content)
	})

	t.Run("should return 404 when note not found", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		noteID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(domain.Note{}, domain.ErrNoteNotFound).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID.Hex())
			handler.GetUserNoteByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID.Hex())
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var got sharedresp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "note not found", got.Message)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		noteID := primitive.NewObjectID()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", handler.GetUserNoteByID)

		uri := fmt.Sprintf("/notes/%s", noteID.Hex())
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var got sharedresp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockRepo := mocks.NewMockNoteRepository(t)
		svc := application.NewNoteService(mockRepo)
		handler := presentation.NewNoteHandler(svc)

		noteID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		mockRepo.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(domain.Note{}, assert.AnError).Once()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID.Hex())
			handler.GetUserNoteByID(w, r.WithContext(ctx))
		})

		uri := fmt.Sprintf("/notes/%s", noteID.Hex())
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, uri, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var got sharedresp.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&got)
		require.NoError(t, err)

		assert.Equal(t, "internal server error", got.Message)
	})
}
