package application

import (
	"context"
	"pinnado/internal/notes/domain"
	"pinnado/pkg/listopts"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) error
	ListByUserID(ctx context.Context, userID primitive.ObjectID, params listopts.ListParams) (listopts.Paginated[domain.Note], error)
}
