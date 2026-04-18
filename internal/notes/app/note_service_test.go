package app_test

import (
	"testing"
	"time"

	"momento/internal/notes/app"
	"momento/internal/notes/domain"
	"momento/internal/notes/mocks"
	"momento/pkg/listopts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewNoteService(t *testing.T) {
	t.Run("should create note service", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		assert.NotNil(t, noteService)
	})
}

func TestNoteService_CreateNote(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	defaultInput := app.NoteInput{
		UserID:  userID,
		Title:   "Title",
		Content: "Note content",
	}

	t.Run("should create note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := noteService.CreateNote(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, domain.NoteTitle("Title"), got.Title)
		assert.Equal(t, domain.NoteContent("Note content"), got.Content)
		assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, time.Second)
		assert.Equal(t, time.UTC, got.CreatedAt.Location())
		assert.Equal(t, time.UTC, got.UpdatedAt.Location())
	})

	t.Run("should return error when title is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.Title = ""

		_, err := noteService.CreateNote(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrTitleEmpty)
	})

	t.Run("should return error when content is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.Content = ""

		_, err := noteService.CreateNote(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid-user-id"

		_, err := noteService.CreateNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when UserID is empty", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = ""

		_, err := noteService.CreateNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when repository Create fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

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

	defaultInput := app.ListNotesInput{
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
		noteService := app.NewNoteService(noteRepoMock)

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
		noteService := app.NewNoteService(noteRepoMock)

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
		noteService := app.NewNoteService(noteRepoMock)

		input := app.ListNotesInput{
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
		noteService := app.NewNoteService(noteRepoMock)

		input := app.ListNotesInput{
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
		noteService := app.NewNoteService(noteRepoMock)

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
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid-user-id"

		_, err := noteService.ListNotes(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when UserID is empty", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = ""

		_, err := noteService.ListNotes(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when repository ListByUserID fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

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

	defaultInput := app.GetUserNoteByIDInput{
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
		noteService := app.NewNoteService(noteRepoMock)

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
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"
		_, err := noteService.GetUserNoteByID(t.Context(), input)
		assert.EqualError(t, err, "invalid ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should throw error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"
		_, err := noteService.GetUserNoteByID(t.Context(), input)
		assert.EqualError(t, err, "invalid user ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should throw error when GetUserNoteByID return note not found", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(domain.Note{}, domain.ErrNoteNotFound)

		_, err := noteService.GetUserNoteByID(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("should throw error when GetUserNoteByID fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(domain.Note{}, assert.AnError)

		_, err := noteService.GetUserNoteByID(t.Context(), defaultInput)

		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestNoteService_UpdateNote(t *testing.T) {
	userID := primitive.NewObjectID()
	noteID := primitive.NewObjectID()

	defaultInput := app.UpdateNoteInput{
		UserID:  userID.Hex(),
		ID:      noteID.Hex(),
		Title:   "Updated title",
		Content: "Updated content",
	}

	now := time.Now().UTC().Add(-time.Hour)
	noteMockDefault := domain.Note{
		ID:        noteID,
		UserID:    userID,
		Title:     "Title",
		Content:   "Initial content",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("should update note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		nMock := noteMockDefault
		nMock.UpdatedAt = time.Now().UTC()

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(nMock, nil).Once()
		noteRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := noteService.UpdateNote(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Equal(t, noteID.Hex(), got.ID)
		assert.Equal(t, userID.Hex(), got.UserID)
		assert.Equal(t, domain.NoteContent("Updated content"), got.Content)
		assert.Equal(t, noteMockDefault.CreatedAt, got.CreatedAt)
		assert.NotEqual(t, noteMockDefault.UpdatedAt, got.UpdatedAt)
	})

	t.Run("should return error when ID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.EqualError(t, err, "invalid ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.EqualError(t, err, "invalid user ID: the provided hex string is not a valid ObjectID")
	})

	t.Run("should return error when title is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.Title = ""

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrTitleEmpty)
	})

	t.Run("should return error when content is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.Content = ""

		_, err := noteService.UpdateNote(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return note not found when repository Update returns ErrNoteNotFound", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(noteMockDefault, nil).Once()
		noteRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.ErrNoteNotFound).Once()

		_, err := noteService.UpdateNote(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("should return error when repository Update fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(noteMockDefault, nil).Once()
		noteRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(assert.AnError).Once()

		_, err := noteService.UpdateNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.Update")
	})

	t.Run("should return error when repository GetByIDAndUserID fails", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, noteID, userID).Return(domain.Note{}, assert.AnError).Once()

		_, err := noteService.UpdateNote(t.Context(), defaultInput)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Contains(t, err.Error(), "s.noteRepository.GetByIDAndUserID")
	})
}

func TestNoteService_DeleteNote(t *testing.T) {
	userID := primitive.NewObjectID()
	noteID := primitive.NewObjectID()

	defaultInput := app.DeleteNoteInput{
		UserID: userID.Hex(),
		ID:     noteID.Hex(),
	}

	t.Run("should delete note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().DeleteByIDAndUserID(mock.Anything, noteID, userID).Return(nil).Once()

		err := noteService.DeleteNote(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when ID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"

		err := noteService.DeleteNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"

		err := noteService.DeleteNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return note not found when repository returns ErrNoteNotFound", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().DeleteByIDAndUserID(mock.Anything, noteID, userID).Return(domain.ErrNoteNotFound).Once()

		err := noteService.DeleteNote(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("should return wrapped error when repository returns generic error", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().DeleteByIDAndUserID(mock.Anything, noteID, userID).Return(assert.AnError).Once()

		err := noteService.DeleteNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.DeleteByIDAndUserID")
	})
}

func TestNoteService_ArchiveNote(t *testing.T) {
	userID := primitive.NewObjectID()
	noteID := primitive.NewObjectID()

	defaultInput := app.ArchiveNoteInput{
		UserID: userID.Hex(),
		ID:     noteID.Hex(),
	}

	t.Run("should archive note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().ArchiveByIDAndUserID(mock.Anything, noteID, userID).Return(nil).Once()

		err := noteService.ArchiveNote(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when ID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"

		err := noteService.ArchiveNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"

		err := noteService.ArchiveNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when note not found", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().ArchiveByIDAndUserID(mock.Anything, noteID, userID).Return(domain.ErrNoteNotFound).Once()

		err := noteService.ArchiveNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("should return wrapped error when repository returns generic error", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().ArchiveByIDAndUserID(mock.Anything, noteID, userID).Return(assert.AnError).Once()

		err := noteService.ArchiveNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.ArchiveByIDAndUserID")
	})
}

func TestNoteService_RestoreNote(t *testing.T) {
	userID := primitive.NewObjectID()
	noteID := primitive.NewObjectID()

	defaultInput := app.RestoreNoteInput{
		UserID: userID.Hex(),
		ID:     noteID.Hex(),
	}

	t.Run("should restore note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().RestoreByIDAndUserID(mock.Anything, noteID, userID).Return(nil).Once()

		err := noteService.RestoreNote(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when ID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.ID = "invalid"

		err := noteService.RestoreNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})

	t.Run("should return error when UserID is invalid", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		input := defaultInput
		input.UserID = "invalid"

		err := noteService.RestoreNote(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error when note not found", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().RestoreByIDAndUserID(mock.Anything, noteID, userID).Return(domain.ErrNoteNotFound).Once()

		err := noteService.RestoreNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("should return wrapped error when repository returns generic error", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := app.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().RestoreByIDAndUserID(mock.Anything, noteID, userID).Return(assert.AnError).Once()

		err := noteService.RestoreNote(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.noteRepository.RestoreByIDAndUserID")
	})
}
