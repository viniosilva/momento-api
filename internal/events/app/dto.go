package app

import (
	"time"

	"momento/internal/events/domain"
	"momento/pkg/listopts"
)

type EventInput struct {
	UserID  string
	Title   string
	Content string
}

type EventOutput struct {
	ID          string
	OwnerUserID string
	Title       domain.EventTitle
	Content     domain.EventContent
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

type ListEventsInput struct {
	UserID     string
	Pagination listopts.PaginationInput
	Sort       listopts.SortInput
}

type ListEventsOutput listopts.Paginated[EventOutput]

type GetUserEventByIDInput struct {
	UserID string
	ID     string
}

type UpdateEventInput struct {
	UserID  string
	ID      string
	Title   *string
	Content *string
}

type DeleteEventInput struct {
	UserID string
	ID     string
}

type ArchiveEventInput struct {
	UserID string
	ID     string
}

type RestoreEventInput struct {
	UserID string
	ID     string
}

type GetUploadURLInput struct {
	UserID      string
	EventID     string
	ContentType string
}

type GetUploadURLOutput struct {
	UploadURL string
	ObjectKey string
}

type ConfirmImageInput struct {
	UserID    string
	EventID   string
	ObjectKey string
}

type ConfirmImageOutput struct {
	Path        domain.ImagePath
	DownloadURL string
}

type DeleteImageInput struct {
	UserID  string
	EventID string
	Path    string
}

type ListImagesInput struct {
	UserID  string
	EventID string
}

type ImageOutput struct {
	Path        string
	DownloadURL string
}

func EventApplicationToOutput(event domain.Event) EventOutput {
	return EventOutput{
		ID:          event.ID,
		OwnerUserID: event.OwnerUserID,
		Title:       event.Title,
		Content:     event.Content,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		ArchivedAt:  event.ArchivedAt,
	}
}
