package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"momento/internal/events/domain"
	"momento/pkg/listopts"
	"momento/pkg/tools"

	"github.com/google/uuid"
)

const (
	PREFIX_PATH_EVENT_IMAGE = "events"
)

type eventService struct {
	eventRepository        EventRepository
	s3Service              S3Service
	imageDownloadURLExpiry time.Duration
}

func NewEventService(eventRepository EventRepository, s3Service S3Service, imageDownloadURLExpiry time.Duration) *eventService {
	return &eventService{
		eventRepository:        eventRepository,
		s3Service:              s3Service,
		imageDownloadURLExpiry: imageDownloadURLExpiry,
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
	event, err := s.eventRepository.GetByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return domain.ErrEventNotFound
		}

		return fmt.Errorf("s.eventRepository.GetByIDAndUserID: %w", err)
	}

	if event.Metadata != nil {
		eventPath := fmt.Sprintf("%s/%s", PREFIX_PATH_EVENT_IMAGE, input.ID)
		if err := s.s3Service.DeleteFolder(ctx, eventPath); err != nil {
			return fmt.Errorf("s.s3Service.DeleteFolder: %w", err)
		}
	}

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

func (s *eventService) GetUploadURL(ctx context.Context, input GetUploadURLInput) (GetUploadURLOutput, error) {
	extension, err := contentTypeToExtension(input.ContentType)
	if err != nil {
		return GetUploadURLOutput{}, err
	}

	objectKey := fmt.Sprintf("%s/%s/%s.%s", PREFIX_PATH_EVENT_IMAGE, input.EventID, uuid.NewString(), extension)

	uploadURL, err := s.s3Service.GetPresignedUploadURL(ctx, objectKey, input.ContentType, time.Hour)
	if err != nil {
		return GetUploadURLOutput{}, fmt.Errorf("s.s3Service.GetPresignedUploadURL: %w", err)
	}

	return GetUploadURLOutput{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	}, nil
}

func (s *eventService) ConfirmImage(ctx context.Context, input ConfirmImageInput) (ConfirmImageOutput, error) {
	path, err := domain.NewImagePath(input.ObjectKey)
	if err != nil {
		return ConfirmImageOutput{}, err
	}

	event, err := s.eventRepository.GetByIDAndUserID(ctx, input.EventID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return ConfirmImageOutput{}, domain.ErrEventNotFound
		}

		return ConfirmImageOutput{}, fmt.Errorf("s.eventRepository.GetByIDAndUserID: %w", err)
	}

	if err := event.AddImage(path); err != nil {
		return ConfirmImageOutput{}, err
	}

	if err := s.eventRepository.AddImage(ctx, input.EventID, input.UserID, path); err != nil {
		return ConfirmImageOutput{}, fmt.Errorf("s.eventRepository.AddImage: %w", err)
	}

	downloadURL, err := s.s3Service.GetPresignedDownloadURL(ctx, input.ObjectKey, s.imageDownloadURLExpiry)
	if err != nil {
		return ConfirmImageOutput{}, fmt.Errorf("s.s3Service.GetPresignedDownloadURL: %w", err)
	}

	return ConfirmImageOutput{
		Path:        path,
		DownloadURL: downloadURL,
	}, nil
}

func (s *eventService) ListImages(ctx context.Context, input ListImagesInput) ([]ImageOutput, error) {
	event, err := s.eventRepository.GetByIDAndUserID(ctx, input.EventID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return nil, domain.ErrEventNotFound
		}

		return nil, fmt.Errorf("s.eventRepository.GetByIDAndUserID: %w", err)
	}

	if event.Metadata == nil || len(event.Metadata.ImagePaths) == 0 {
		return make([]ImageOutput, 0), nil
	}

	outputs := make([]ImageOutput, 0, len(event.Metadata.ImagePaths))
	for _, imagePath := range event.Metadata.ImagePaths {
		downloadURL, err := s.s3Service.GetPresignedDownloadURL(ctx, string(imagePath), s.imageDownloadURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("s.s3Service.GetPresignedDownloadURL: %w", err)
		}

		outputs = append(outputs, ImageOutput{
			Path:        string(imagePath),
			DownloadURL: downloadURL,
		})
	}

	return outputs, nil
}

func (s *eventService) DeleteImage(ctx context.Context, input DeleteImageInput) error {
	path, err := domain.NewImagePath(input.Path)
	if err != nil {
		return fmt.Errorf("domain.NewImagePath: %w", err)
	}

	if err := s.eventRepository.RemoveImage(ctx, input.EventID, input.UserID, path); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return domain.ErrEventNotFound
		}

		return fmt.Errorf("s.eventRepository.RemoveImage: %w", err)
	}

	if err := s.s3Service.DeleteObject(ctx, input.Path); err != nil {
		return fmt.Errorf("s.s3Service.DeleteObject: %w", err)
	}

	return nil
}

func contentTypeToExtension(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "image/webp":
		return "webp", nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
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
