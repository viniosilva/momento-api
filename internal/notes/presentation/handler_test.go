package presentation_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	"pinnado/internal/notes/presentation"
	sharedresp "pinnado/internal/shared/presentation/response"
	"pinnado/mocks"
	"pinnado/pkg/nethttp"
	"pinnado/pkg/pagination"
)

var mapErrorToHTTPStatus = presentation.MapErrorToHTTPStatus

func TestNewNoteHandler(t *testing.T) {
	t.Run("should create note handler", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		assert.NotNil(t, handler)
	})
}

func TestNoteHandler_CreateNote(t *testing.T) {
	validUserID := "507f1f77bcf86cd799439011"

	t.Run("should return 201 when note is created successfully", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		output := application.NoteOutput{
			ID:        "note123",
			UserID:    validUserID,
			Content:   domain.NoteContent("Valid note content"),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		mockService.EXPECT().CreateNote(mock.Anything, mock.Anything).
			Return(output, nil).
			Once()

		reqBody := map[string]any{
			"content": "Valid note content",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, presentation.NoteResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "note123", got.ID)
		assert.Equal(t, validUserID, got.UserID)
		assert.Equal(t, "Valid note content", got.Content)
	})

	t.Run("should return 400 when content is empty", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		mockService.EXPECT().CreateNote(mock.Anything, mock.Anything).
			Return(application.NoteOutput{}, domain.ErrContentEmpty).
			Once()

		reqBody := map[string]any{
			"content": "",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, got.Message, "content cannot be empty")
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

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
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		mockService.EXPECT().CreateNote(mock.Anything, mock.Anything).
			Return(application.NoteOutput{}, assert.AnError).
			Once()

		reqBody := map[string]any{
			"content": "Valid note content",
		}

		resp, got, err := nethttp.RequestWithResponse[map[string]any, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", reqBody, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.CreateNote(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})

	t.Run("should return 400 when request body is invalid JSON", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		resp, got, err := nethttp.RequestWithResponse[string, sharedresp.ErrorResponse](
			t.Context(), http.MethodPost, "/notes", "invalid json", func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
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
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		now := time.Now().UTC()
		output := application.ListNotesOutput{
			Data: []application.NoteOutput{
				{
					ID:        "note1",
					UserID:    validUserID,
					Content:   domain.NoteContent("Note 1 content"),
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "note2",
					UserID:    validUserID,
					Content:   domain.NoteContent("Note 2 content"),
					CreatedAt: now.Add(-time.Hour),
					UpdatedAt: now.Add(-time.Hour),
				},
			},
			Pagination: pagination.PaginationOutput{
				TotalCount: 2,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			},
		}

		mockService.EXPECT().ListNotes(mock.Anything, mock.Anything).
			Return(output, nil).
			Once()

		resp, got, err := nethttp.RequestWithResponse[any, pagination.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes?page=1&page_size=10", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, got.Data, 2)
		assert.Equal(t, int64(2), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 1, got.Pagination.TotalPages)
		assert.Equal(t, "note1", got.Data[0].ID)
		assert.Equal(t, "Note 1 content", got.Data[0].Content)
	})

	t.Run("should return 200 with empty list when no notes found", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		output := application.ListNotesOutput{
			Data: []application.NoteOutput{},
			Pagination: pagination.PaginationOutput{
				TotalCount: 0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
			},
		}

		mockService.EXPECT().ListNotes(mock.Anything, mock.Anything).
			Return(output, nil).
			Once()

		resp, got, err := nethttp.RequestWithResponse[any, pagination.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Empty(t, got.Data)
		assert.Equal(t, int64(0), got.Pagination.TotalCount)
	})

	t.Run("should use default values when query parameters are missing", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		output := application.ListNotesOutput{
			Data: []application.NoteOutput{},
			Pagination: pagination.PaginationOutput{
				TotalCount: 0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
			},
		}

		mockService.EXPECT().ListNotes(mock.Anything, mock.Anything).Return(output, nil).Once()

		resp, _, err := nethttp.RequestWithResponse[any, pagination.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should parse query parameters correctly", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		output := application.ListNotesOutput{
			Data: []application.NoteOutput{},
			Pagination: pagination.PaginationOutput{
				TotalCount: 0,
				Page:       2,
				PageSize:   10,
				TotalPages: 0,
			},
		}

		mockService.EXPECT().ListNotes(mock.Anything, mock.Anything).Return(output, nil).Once()

		resp, _, err := nethttp.RequestWithResponse[any, pagination.PaginatedResponse[presentation.NoteResponse]](
			t.Context(), http.MethodGet, "/notes?page=2&page_size=10&sort_by=updated_at&sort_order=asc", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 401 when UserID is missing from context", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		resp, got, err := nethttp.RequestWithResponse[any, sharedresp.ErrorResponse](
			t.Context(), http.MethodGet, "/notes", nil, handler.ListNotes)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, "unauthorized", got.Message)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockService := mocks.NewMockNoteService(t)
		handler := presentation.NewNoteHandler(mockService)

		mockService.EXPECT().ListNotes(mock.Anything, mock.Anything).
			Return(application.ListNotesOutput{}, assert.AnError).
			Once()

		resp, got, err := nethttp.RequestWithResponse[any, sharedresp.ErrorResponse](
			t.Context(), http.MethodGet, "/notes", nil, func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), nethttp.ContextKeyUserID, validUserID)
				handler.ListNotes(w, r.WithContext(ctx))
			})
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "internal server error", got.Message)
	})
}
