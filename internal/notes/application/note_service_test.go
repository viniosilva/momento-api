package application_test

import (
	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	"pinnado/internal/notes/mocks"
	"pinnado/pkg/listopts"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewNoteService(t *testing.T) {
	t.Run("should create note service", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		assert.NotNil(t, noteService)
	})
}

func TestNoteService_CreateNote(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	defaultInput := application.NoteInput{
		UserID:  userID,
		Content: "Valid note content",
	}

	t.Run("should create note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := noteService.CreateNote(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, domain.NoteContent("Valid note content"), got.Content)
		assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, time.Second)
		assert.Equal(t, time.UTC, got.CreatedAt.Location())
		assert.Equal(t, time.UTC, got.UpdatedAt.Location())
	})

	t.Run("should return error when content is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.Content = ""

		_, err := noteService.CreateNote(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid-user-id"

		_, err := noteService.CreateNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when UserID is empty", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = ""

		_, err := noteService.CreateNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when repository Create fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		_, err := noteService.CreateNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.Create")
	})
}

func TestNoteService_ListNotes(t *testing.T) {
	userID := primitive.NewObjectID()

	defaultInput := application.ListNotesInput{
		UserID: userID.Hex(),
		Pagination: listopts.PaginationInput{
			Page:     1,
			PageSize: 10,
		},
		Sort: listopts.SortInput{
			Field: "created_at",
			Order: "desc",
		},
	}

	t.Run("should list notes successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		expectedNotes := []domain.Note{
			{
				ID:        primitive.NewObjectID(),
				UserID:    userID,
				Content:   "Note 1",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			{
				ID:        primitive.NewObjectID(),
				UserID:    userID,
				Content:   "Note 2",
				CreatedAt: time.Now().UTC().Add(-time.Hour),
				UpdatedAt: time.Now().UTC().Add(-time.Hour),
			},
		}

		noteRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Note]{
				Data: expectedNotes,
				Pagination: listopts.PaginationOutput{
					TotalCount: 2,
					Page:       1,
					PageSize:   10,
					TotalPages: 1,
				},
			}, nil).
			Once()

		got, err := noteService.ListNotes(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Len(t, got.Data, 2)
		assert.Equal(t, int64(2), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 1, got.Pagination.TotalPages)
		assert.Equal(t, expectedNotes[0].ID.Hex(), got.Data[0].ID)
		assert.Equal(t, expectedNotes[1].ID.Hex(), got.Data[1].ID)
	})

	t.Run("should return empty list when no notes found", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Note]{
				Data: []domain.Note{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				},
			}, nil).
			Once()

		got, err := noteService.ListNotes(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Empty(t, got.Data)
		assert.Equal(t, int64(0), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 0, got.Pagination.TotalPages)
	})

	t.Run("should calculate total pages correctly", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := application.ListNotesInput{
			UserID: userID.Hex(),
			Pagination: listopts.PaginationInput{
				Page:     1,
				PageSize: 10,
			},
			Sort: listopts.SortInput{
				Field: "created_at",
				Order: "desc",
			},
		}

		noteRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Note]{
				Data: []domain.Note{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 25,
					Page:       1,
					PageSize:   10,
					TotalPages: 3,
				},
			}, nil).
			Once()

		got, err := noteService.ListNotes(t.Context(), input)
		require.NoError(t, err)

		assert.Equal(t, int64(25), got.Pagination.TotalCount)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 3, got.Pagination.TotalPages)
	})

	t.Run("should apply default pagination when invalid values provided", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := application.ListNotesInput{
			UserID: userID.Hex(),
			Pagination: listopts.PaginationInput{
				Page:     0,
				PageSize: 0,
			},
			Sort: listopts.SortInput{
				Field: "created_at",
				Order: "desc",
			},
		}

		noteRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Note]{
				Data: []domain.Note{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				},
			}, nil).
			Once()

		got, err := noteService.ListNotes(t.Context(), input)
		require.NoError(t, err)

		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
	})

	t.Run("should apply default sort when invalid values provided", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.Sort = listopts.SortInput{
			Field: "invalid_field",
			Order: "invalid_order",
		}

		noteRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Note]{
				Data: []domain.Note{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				},
			}, nil).
			Once()

		_, err := noteService.ListNotes(t.Context(), input)
		require.NoError(t, err)
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid-user-id"

		_, err := noteService.ListNotes(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when UserID is empty", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = ""

		_, err := noteService.ListNotes(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when repository ListByUserID fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Note]{}, assert.AnError).
			Once()

		_, err := noteService.ListNotes(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.ListByUserID")
	})
}

func TestNoteService_GetUserNoteByID(t *testing.T) {
	userID := primitive.NewObjectID()
	noteID := primitive.NewObjectID()

	defaultInput := application.GetUserNoteByIDInput{
		UserID: userID.Hex(),
		ID:     userID.Hex(),
	}

	now := time.Now()
	noteMock := domain.Note{
		ID:        noteID,
		UserID:    userID,
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("should get user's note by ID", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(noteMock, nil)

		got, err := noteService.GetUserNoteByID(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Equal(t, noteMock.ID.Hex(), got.ID)
		assert.Equal(t, noteMock.UserID.Hex(), got.UserID)
		assert.Equal(t, noteMock.Content, got.Content)
		assert.Equal(t, noteMock.CreatedAt, got.CreatedAt)
		assert.Equal(t, noteMock.UpdatedAt, got.UpdatedAt)
	})

	t.Run("should throw error when ID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"
		_, err := noteService.GetUserNoteByID(t.Context(), input)
		assert.EqualError(t, err, "invalid ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should throw error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"
		_, err := noteService.GetUserNoteByID(t.Context(), input)
		assert.EqualError(t, err, "invalid user ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should throw error when GetUserNoteByID fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(domain.Note{}, assert.AnError)

		_, err := noteService.GetUserNoteByID(t.Context(), defaultInput)

		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestNoteService_UpdateNote(t *testing.T) {
	userID := primitive.NewObjectID()
	noteID := primitive.NewObjectID()

	defaultInput := application.UpdateNoteInput{
		UserID:  userID.Hex(),
		ID:      noteID.Hex(),
		Content: "Updated content",
	}

	now := time.Now().UTC().Add(-time.Hour)
	noteMock := domain.Note{
		ID:        noteID,
		UserID:    userID,
		Content:   "Initial content",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("should update note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(noteMock, nil).Once()
		noteRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := noteService.UpdateNote(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Equal(t, noteID.Hex(), got.ID)
		assert.Equal(t, userID.Hex(), got.UserID)
		assert.Equal(t, domain.NoteContent("Updated content"), got.Content)
		assert.Equal(t, noteMock.CreatedAt, got.CreatedAt)
		assert.NotEqual(t, noteMock.UpdatedAt, got.UpdatedAt)
		assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, time.Second)
	})

	t.Run("should return error when ID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.EqualError(t, err, "invalid ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.EqualError(t, err, "invalid user ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should return error when content is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(noteMock, nil).Once()

		input := defaultInput
		input.Content = ""

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return note not found when repository Update returns ErrNoteNotFound", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(noteMock, nil).Once()
		noteRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.ErrNoteNotFound).Once()

		_, err := noteService.UpdateNote(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("should return error when repository Update fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(noteMock, nil).Once()
		noteRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(assert.AnError).Once()

		_, err := noteService.UpdateNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.Update")
	})

	t.Run("should return error when repository GetByIDAndUserID fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(domain.Note{}, assert.AnError).Once()

		_, err := noteService.UpdateNote(t.Context(), defaultInput)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Contains(t, err.Error(), "s.noteRepository.GetByIDAndUserID")
	})
}
