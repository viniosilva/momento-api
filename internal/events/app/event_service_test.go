package app_test

import (
	"testing"
	"time"

	"momento/internal/events/app"
	"momento/internal/events/domain"
	"momento/internal/events/mocks"
	"momento/pkg/listopts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewEventService(t *testing.T) {
	t.Run("should create event service", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		assert.NotNil(t, eventService)
	})
}

func TestEventService_CreateEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	defaultInput := app.EventInput{
		UserID:  userID,
		Title:   "Title",
		Content: "Event content",
	}

	t.Run("should create event successfully", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := eventService.CreateEvent(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, domain.EventTitle("Title"), got.Title)
		assert.Equal(t, domain.EventContent("Event content"), got.Content)
		assert.WithinDuration(t, time.Now().UTC(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, time.Second)
		assert.Equal(t, time.UTC, got.CreatedAt.Location())
		assert.Equal(t, time.UTC, got.UpdatedAt.Location())
	})

	t.Run("should return error when title is invalid", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := defaultInput
		input.Title = ""

		_, err := eventService.CreateEvent(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrTitleEmpty)
	})

	t.Run("should return error when content is invalid", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := defaultInput
		input.Content = ""

		_, err := eventService.CreateEvent(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return error when repository Create fails", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().Create(mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		_, err := eventService.CreateEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.eventRepository.Create")
	})
}

func TestEventService_ListEvents(t *testing.T) {
	userID := primitive.NewObjectID().Hex()

	defaultInput := app.ListEventsInput{
		UserID: userID,
		Pagination: listopts.PaginationInput{
			Page:     1,
			PageSize: 10,
		},
		Sort: listopts.SortInput{
			Field: "created_at",
			Order: "desc",
		},
	}

	t.Run("should list events successfully", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		expectedEvents := []domain.Event{
			{
				ID:          primitive.NewObjectID().Hex(),
				OwnerUserID: userID,
				Content:     "Event 1",
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			},
			{
				ID:          primitive.NewObjectID().Hex(),
				OwnerUserID: userID,
				Content:     "Event 2",
				CreatedAt:   time.Now().UTC().Add(-time.Hour),
				UpdatedAt:   time.Now().UTC().Add(-time.Hour),
			},
		}

		eventRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Event]{
				Data: expectedEvents,
				Pagination: listopts.PaginationOutput{
					TotalCount: 2,
					Page:       1,
					PageSize:   10,
					TotalPages: 1,
				},
			}, nil).
			Once()

		got, err := eventService.ListEvents(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Len(t, got.Data, 2)
		assert.Equal(t, int64(2), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 1, got.Pagination.TotalPages)
		assert.Equal(t, expectedEvents[0].ID, got.Data[0].ID)
		assert.Equal(t, expectedEvents[1].ID, got.Data[1].ID)
	})

	t.Run("should return empty list when no events found", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Event]{
				Data: []domain.Event{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				},
			}, nil).
			Once()

		got, err := eventService.ListEvents(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Empty(t, got.Data)
		assert.Equal(t, int64(0), got.Pagination.TotalCount)
		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 0, got.Pagination.TotalPages)
	})

	t.Run("should calculate total pages correctly", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := app.ListEventsInput{
			UserID: userID,
			Pagination: listopts.PaginationInput{
				Page:     1,
				PageSize: 10,
			},
			Sort: listopts.SortInput{
				Field: "created_at",
				Order: "desc",
			},
		}

		eventRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Event]{
				Data: []domain.Event{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 25,
					Page:       1,
					PageSize:   10,
					TotalPages: 3,
				},
			}, nil).
			Once()

		got, err := eventService.ListEvents(t.Context(), input)
		require.NoError(t, err)

		assert.Equal(t, int64(25), got.Pagination.TotalCount)
		assert.Equal(t, 10, got.Pagination.PageSize)
		assert.Equal(t, 3, got.Pagination.TotalPages)
	})

	t.Run("should apply default pagination when invalid values provided", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := app.ListEventsInput{
			UserID: userID,
			Pagination: listopts.PaginationInput{
				Page:     0,
				PageSize: 0,
			},
			Sort: listopts.SortInput{
				Field: "created_at",
				Order: "desc",
			},
		}

		eventRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Event]{
				Data: []domain.Event{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				},
			}, nil).
			Once()

		got, err := eventService.ListEvents(t.Context(), input)
		require.NoError(t, err)

		assert.Equal(t, 1, got.Pagination.Page)
		assert.Equal(t, 10, got.Pagination.PageSize)
	})

	t.Run("should apply default sort when invalid values provided", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := defaultInput
		input.Sort = listopts.SortInput{
			Field: "invalid_field",
			Order: "invalid_order",
		}

		eventRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Event]{
				Data: []domain.Event{},
				Pagination: listopts.PaginationOutput{
					TotalCount: 0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				},
			}, nil).
			Once()

		_, err := eventService.ListEvents(t.Context(), input)
		require.NoError(t, err)
	})

	t.Run("should return error when repository ListByUserID fails", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().
			ListByUserID(mock.Anything, userID, mock.Anything).
			Return(listopts.Paginated[domain.Event]{}, assert.AnError).
			Once()

		_, err := eventService.ListEvents(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.eventRepository.ListByUserID")
	})
}

func TestEventService_GetUserEventByID(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	eventID := primitive.NewObjectID().Hex()

	defaultInput := app.GetUserEventByIDInput{
		UserID: userID,
		ID:     eventID,
	}

	now := time.Now()
	eventMock := domain.Event{
		ID:          eventID,
		OwnerUserID: userID,
		Content:     "Content",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("should get user's event by ID", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(eventMock, nil)

		got, err := eventService.GetUserEventByID(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Equal(t, eventMock.ID, got.ID)
		assert.Equal(t, eventMock.OwnerUserID, got.UserID)
		assert.Equal(t, eventMock.Content, got.Content)
		assert.Equal(t, eventMock.CreatedAt, got.CreatedAt)
		assert.Equal(t, eventMock.UpdatedAt, got.UpdatedAt)
	})

	t.Run("should throw error when GetUserEventByID return event not found", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(domain.Event{}, domain.ErrEventNotFound)

		_, err := eventService.GetUserEventByID(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("should throw error when GetUserEventByID fails", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, mock.Anything, mock.Anything).Return(domain.Event{}, assert.AnError)

		_, err := eventService.GetUserEventByID(t.Context(), defaultInput)

		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEventService_UpdateEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	eventID := primitive.NewObjectID().Hex()

	defaultInput := app.UpdateEventInput{
		UserID:  userID,
		ID:      eventID,
		Title:   new("Updated title"),
		Content: new("Updated content"),
	}

	now := time.Now().UTC().Add(-time.Hour)
	eventMockDefault := domain.Event{
		ID:          eventID,
		OwnerUserID: userID,
		Title:       "Title",
		Content:     "Initial content",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("should update event successfully", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eMock := eventMockDefault
		eMock.UpdatedAt = time.Now().UTC()

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(eMock, nil).Once()
		eventRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		got, err := eventService.UpdateEvent(t.Context(), defaultInput)
		require.NoError(t, err)

		assert.Equal(t, eventID, got.ID)
		assert.Equal(t, userID, got.UserID)
		assert.Equal(t, domain.EventTitle("Updated title"), got.Title)
		assert.Equal(t, domain.EventContent("Updated content"), got.Content)
		assert.Equal(t, eventMockDefault.CreatedAt, got.CreatedAt)
		assert.NotEqual(t, eventMockDefault.UpdatedAt, got.UpdatedAt)
	})

	t.Run("should update only title when content is nil", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(eventMockDefault, nil).Once()
		eventRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		input := app.UpdateEventInput{
			UserID:  userID,
			ID:      eventID,
			Title:   new("New title"),
			Content: nil,
		}

		got, err := eventService.UpdateEvent(t.Context(), input)
		require.NoError(t, err)

		assert.Equal(t, domain.EventTitle("New title"), got.Title)
		assert.Equal(t, eventMockDefault.Content, got.Content)
	})

	t.Run("should update only content when title is nil", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(eventMockDefault, nil).Once()
		eventRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		input := app.UpdateEventInput{
			UserID:  userID,
			ID:      eventID,
			Title:   nil,
			Content: new("New content"),
		}

		got, err := eventService.UpdateEvent(t.Context(), input)
		require.NoError(t, err)

		assert.Equal(t, eventMockDefault.Title, got.Title)
		assert.Equal(t, domain.EventContent("New content"), got.Content)
	})

	t.Run("should return error when title is invalid", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := defaultInput
		input.Title = new("")

		_, err := eventService.UpdateEvent(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrTitleEmpty)
	})

	t.Run("should return error when content is invalid", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		input := defaultInput
		input.Content = new("")

		_, err := eventService.UpdateEvent(t.Context(), input)

		assert.ErrorIs(t, err, domain.ErrContentEmpty)
	})

	t.Run("should return event not found when repository Update returns ErrEventNotFound", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(eventMockDefault, nil).Once()
		eventRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.ErrEventNotFound).Once()

		_, err := eventService.UpdateEvent(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("should return error when repository Update fails", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(eventMockDefault, nil).Once()
		eventRepoMock.EXPECT().Update(mock.Anything, mock.Anything).Return(assert.AnError).Once()

		_, err := eventService.UpdateEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.eventRepository.Update")
	})

	t.Run("should return error when repository GetByIDAndUserID fails", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().GetByIDAndUserID(mock.Anything, eventID, userID).Return(domain.Event{}, assert.AnError).Once()

		_, err := eventService.UpdateEvent(t.Context(), defaultInput)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Contains(t, err.Error(), "s.eventRepository.GetByIDAndUserID")
	})
}

func TestEventService_DeleteEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	eventID := primitive.NewObjectID().Hex()

	defaultInput := app.DeleteEventInput{
		UserID: userID,
		ID:     eventID,
	}

	t.Run("should delete event successfully", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().DeleteByIDAndUserID(mock.Anything, eventID, userID).Return(nil).Once()

		err := eventService.DeleteEvent(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return event not found when repository returns ErrEventNotFound", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().DeleteByIDAndUserID(mock.Anything, eventID, userID).Return(domain.ErrEventNotFound).Once()

		err := eventService.DeleteEvent(t.Context(), defaultInput)

		assert.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("should return wrapped error when repository returns generic error", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().DeleteByIDAndUserID(mock.Anything, eventID, userID).Return(assert.AnError).Once()

		err := eventService.DeleteEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.eventRepository.DeleteByIDAndUserID")
	})
}

func TestEventService_ArchiveEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	eventID := primitive.NewObjectID().Hex()

	defaultInput := app.ArchiveEventInput{
		UserID: userID,
		ID:     eventID,
	}

	t.Run("should archive event successfully", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().ArchiveByIDAndUserID(mock.Anything, eventID, userID).Return(nil).Once()

		err := eventService.ArchiveEvent(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when event not found", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().ArchiveByIDAndUserID(mock.Anything, eventID, userID).Return(domain.ErrEventNotFound).Once()

		err := eventService.ArchiveEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("should return wrapped error when repository returns generic error", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().ArchiveByIDAndUserID(mock.Anything, eventID, userID).Return(assert.AnError).Once()

		err := eventService.ArchiveEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.eventRepository.ArchiveByIDAndUserID")
	})
}

func TestEventService_RestoreEvent(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	eventID := primitive.NewObjectID().Hex()

	defaultInput := app.RestoreEventInput{
		UserID: userID,
		ID:     eventID,
	}

	t.Run("should restore event successfully", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().RestoreByIDAndUserID(mock.Anything, eventID, userID).Return(nil).Once()

		err := eventService.RestoreEvent(t.Context(), defaultInput)
		require.NoError(t, err)
	})

	t.Run("should return error when event not found", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().RestoreByIDAndUserID(mock.Anything, eventID, userID).Return(domain.ErrEventNotFound).Once()

		err := eventService.RestoreEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("should return wrapped error when repository returns generic error", func(t *testing.T) {
		eventRepoMock := mocks.NewMockEventRepository(t)
		eventService := app.NewEventService(eventRepoMock)

		eventRepoMock.EXPECT().RestoreByIDAndUserID(mock.Anything, eventID, userID).Return(assert.AnError).Once()

		err := eventService.RestoreEvent(t.Context(), defaultInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s.eventRepository.RestoreByIDAndUserID")
	})
}
