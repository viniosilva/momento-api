package presentation

import (
	"context"

	"pinnado/internal/auth/infrastructure"
	"pinnado/internal/notes/application"
)

type NoteService interface {
	CreateNote(ctx context.Context, input application.NoteInput) (application.NoteOutput, error)
	ListNotes(ctx context.Context, input application.ListNotesInput) (application.ListNotesOutput, error)
	GetUserNoteByID(ctx context.Context, input application.GetUserNoteByIDInput) (application.NoteOutput, error)
	UpdateNote(ctx context.Context, input application.UpdateNoteInput) (application.NoteOutput, error)
	DeleteNote(ctx context.Context, input application.DeleteNoteInput) error
	ArchiveNote(ctx context.Context, input application.ArchiveNoteInput) error
	RestoreNote(ctx context.Context, input application.RestoreNoteInput) error
}

type JWTService interface {
	Validate(tokenString string) (*infrastructure.Claims, error)
}
