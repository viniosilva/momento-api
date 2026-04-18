package app

import (
	"context"

	"momento/internal/notes/domain"
	"momento/pkg/listopts"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) error
	ListByUserID(ctx context.Context, userID string, params listopts.ListParams) (listopts.Paginated[domain.Note], error)
	GetByIDAndUserID(ctx context.Context, id, userID string) (domain.Note, error)
	Update(ctx context.Context, note domain.Note) error
	DeleteByIDAndUserID(ctx context.Context, id, userID string) error
	ArchiveByIDAndUserID(ctx context.Context, id, userID string) error
	RestoreByIDAndUserID(ctx context.Context, id, userID string) error
}
