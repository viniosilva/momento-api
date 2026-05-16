package app

import (
	"context"
	"time"

	"momento/internal/events/domain"
	"momento/pkg/listopts"
)

type EventRepository interface {
	Create(ctx context.Context, event domain.Event) error
	ListByUserID(ctx context.Context, userID string, params listopts.ListParams) (listopts.Paginated[domain.Event], error)
	GetByIDAndUserID(ctx context.Context, id, userID string) (domain.Event, error)
	Update(ctx context.Context, event domain.Event) error
	DeleteByIDAndUserID(ctx context.Context, id, userID string) error
	ArchiveByIDAndUserID(ctx context.Context, id, userID string) error
	RestoreByIDAndUserID(ctx context.Context, id, userID string) error
	AddImage(ctx context.Context, eventID, userID string, path domain.ImagePath) error
	RemoveImage(ctx context.Context, eventID, userID string, path domain.ImagePath) error
}

type S3Service interface {
	GetPresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error)
	GetPresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}
