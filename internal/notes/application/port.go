package application

import (
	"context"
	"pinnado/internal/notes/domain"
	shareddto "pinnado/internal/shared/application/dto"
	sharedinfra "pinnado/internal/shared/infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) error
	ListByUserID(ctx context.Context, userID primitive.ObjectID, params sharedinfra.ListParams) (shareddto.Paginated[domain.Note], error)
}
