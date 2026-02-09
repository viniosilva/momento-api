package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pinnado/internal/notes/application"
	"pinnado/internal/notes/domain"
	"pinnado/mocks"

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
	validUserID := primitive.NewObjectID().Hex()

	defaultInput := application.NoteInput{
		UserID:  validUserID,
		Content: "Valid note content",
	}

	t.Run("should create note successfully", func(t *testing.T) {
		noteRepoMock := mocks.NewMockNoteRepository(t)
		noteService := application.NewNoteService(noteRepoMock)

		noteRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := noteService.CreateNote(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, validUserID, got.UserID)
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
