package ports

import (
	"context"

	"momento/internal/events/app"
	appjwt "momento/pkg/jwt"
)

type EventService interface {
	CreateEvent(ctx context.Context, input app.EventInput) (app.EventOutput, error)
	ListEvents(ctx context.Context, input app.ListEventsInput) (app.ListEventsOutput, error)
	GetUserEventByID(ctx context.Context, input app.GetUserEventByIDInput) (app.EventOutput, error)
	UpdateEvent(ctx context.Context, input app.UpdateEventInput) (app.EventOutput, error)
	DeleteEvent(ctx context.Context, input app.DeleteEventInput) error
	ArchiveEvent(ctx context.Context, input app.ArchiveEventInput) error
	RestoreEvent(ctx context.Context, input app.RestoreEventInput) error
	GetUploadURL(ctx context.Context, input app.GetUploadURLInput) (app.GetUploadURLOutput, error)
	ConfirmImage(ctx context.Context, input app.ConfirmImageInput) (app.ConfirmImageOutput, error)
	ListImages(ctx context.Context, input app.ListImagesInput) ([]app.ImageOutput, error)
	DeleteImage(ctx context.Context, input app.DeleteImageInput) error
}

type JWTService interface {
	Validate(tokenString string) (appjwt.UserClaims, error)
}
