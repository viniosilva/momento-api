package app

import (
	"context"
	"errors"
	"fmt"

	"momento/internal/events/domain"
	"momento/pkg/listopts"
	"momento/pkg/tools"
)

type eventService struct {
	eventRepository EventRepository
}

func NewEventService(eventRepository EventRepository) *eventService {
	return &eventService{
		eventRepository: eventRepository,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, input EventInput) (EventOutput, error) {
	title, err := domain.NewEventTitle(input.Title)
	if err != nil {
		return EventOutput{}, err
	}

	content, err := domain.NewEventContent(input.Content)
	if err != nil {
		return EventOutput{}, err
	}

	event := domain.NewEvent(input.UserID, title, content)

	if err := s.eventRepository.Create(ctx, event); err != nil {
		return EventOutput{}, fmt.Errorf("s.eventRepository.Create: %w", err)
	}

	return EventApplicationToOutput(event), nil
}

func (s *eventService) ListEvents(ctx context.Context, input ListEventsInput) (ListEventsOutput, error) {
	params := listopts.ListParams{
		Pagination: input.Pagination,
		Sort:       input.Sort,
	}
	paginatedEvents, err := s.eventRepository.ListByUserID(ctx, input.UserID, params)
	if err != nil {
		return ListEventsOutput{}, fmt.Errorf("s.eventRepository.ListByUserID: %w", err)
	}

	eventOutputs := make([]EventOutput, len(paginatedEvents.Data))
	for i, event := range paginatedEvents.Data {
		eventOutputs[i] = EventApplicationToOutput(event)
	}

	return ListEventsOutput{
		Data:       eventOutputs,
		Pagination: paginatedEvents.Pagination,
	}, nil
}

func (s *eventService) GetUserEventByID(ctx context.Context, input GetUserEventByIDInput) (EventOutput, error) {
	event, err := s.eventRepository.GetByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return EventOutput{}, domain.ErrEventNotFound
		}

		return EventOutput{}, fmt.Errorf("s.eventRepository.GetByIDAndUserID: %w", err)
	}

	return EventApplicationToOutput(event), nil
}

func (s *eventService) UpdateEvent(ctx context.Context, input UpdateEventInput) (EventOutput, error) {
	var newTitle *domain.EventTitle
	if input.Title != nil {
		title, err := domain.NewEventTitle(*input.Title)
		if err != nil {
			return EventOutput{}, err
		}
		newTitle = &title
	}

	var newContent *domain.EventContent
	if input.Content != nil {
		content, err := domain.NewEventContent(*input.Content)
		if err != nil {
			return EventOutput{}, err
		}
		newContent = &content
	}

	event, err := s.eventRepository.GetByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return EventOutput{}, fmt.Errorf("s.eventRepository.GetByIDAndUserID: %w", err)
	}

	event.Update(
		tools.ValueOrDefault(newTitle, event.Title),
		tools.ValueOrDefault(newContent, event.Content),
	)

	if err := s.eventRepository.Update(ctx, event); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return EventOutput{}, domain.ErrEventNotFound
		}

		return EventOutput{}, fmt.Errorf("s.eventRepository.Update: %w", err)
	}

	return EventApplicationToOutput(event), nil
}

func (s *eventService) DeleteEvent(ctx context.Context, input DeleteEventInput) error {
	if err := s.eventRepository.DeleteByIDAndUserID(ctx, input.ID, input.UserID); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return domain.ErrEventNotFound
		}

		return fmt.Errorf("s.eventRepository.DeleteByIDAndUserID: %w", err)
	}

	return nil
}

func (s *eventService) ArchiveEvent(ctx context.Context, input ArchiveEventInput) error {
	err := s.eventRepository.ArchiveByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return domain.ErrEventNotFound
		}

		return fmt.Errorf("s.eventRepository.ArchiveByIDAndUserID: %w", err)
	}

	return nil
}

func (s *eventService) RestoreEvent(ctx context.Context, input RestoreEventInput) error {
	err := s.eventRepository.RestoreByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return domain.ErrEventNotFound
		}

		return fmt.Errorf("s.eventRepository.RestoreByIDAndUserID: %w", err)
	}

	return nil
}