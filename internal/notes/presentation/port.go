package presentation

import (
	"context"

	"pinnado/internal/notes/application"
)

type NoteService interface {
	CreateNote(ctx context.Context, input application.NoteInput) (application.NoteOutput, error)
	ListNotes(ctx context.Context, input application.ListNotesInput) (application.ListNotesOutput, error)
}
