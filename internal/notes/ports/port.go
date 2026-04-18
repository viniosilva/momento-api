package ports

import (
	"context"

	"momento/internal/notes/app"
	appjwt "momento/pkg/jwt"
)

type NoteService interface {
	CreateNote(ctx context.Context, input app.NoteInput) (app.NoteOutput, error)
	ListNotes(ctx context.Context, input app.ListNotesInput) (app.ListNotesOutput, error)
	GetUserNoteByID(ctx context.Context, input app.GetUserNoteByIDInput) (app.NoteOutput, error)
	UpdateNote(ctx context.Context, input app.UpdateNoteInput) (app.NoteOutput, error)
	DeleteNote(ctx context.Context, input app.DeleteNoteInput) error
	ArchiveNote(ctx context.Context, input app.ArchiveNoteInput) error
	RestoreNote(ctx context.Context, input app.RestoreNoteInput) error
}

type JWTService interface {
	Validate(tokenString string) (appjwt.UserClaims, error)
}
