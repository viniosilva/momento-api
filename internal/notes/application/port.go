package application

import (
	"context"
	"pinnado/internal/notes/domain"
	sharedinfra "pinnado/internal/shared/infrastructure"
	"pinnado/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) error
	ListByUserID(ctx context.Context, userID primitive.ObjectID, params sharedinfra.ListParams) (pagination.Paginated[domain.Note], error)
}
