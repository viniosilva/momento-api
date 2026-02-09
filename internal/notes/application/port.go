package application

import (
	"context"

	"pinnado/internal/notes/domain"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) error
}
