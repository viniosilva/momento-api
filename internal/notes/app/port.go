package app

import (
	"context"

	"momento/internal/notes/domain"
	"momento/pkg/listopts"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) error
	ListByUserID(ctx context.Context, userID primitive.ObjectID, params listopts.ListParams) (listopts.Paginated[domain.Note], error)
	GetByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) (domain.Note, error)
	Update(ctx context.Context, note domain.Note) error
	DeleteByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) error
	ArchiveByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) error
	RestoreByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) error
}
