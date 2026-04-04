package application

import (
	"pinnado/internal/notes/domain"
	"pinnado/pkg/listopts"
	"time"
)

type NoteInput struct {
	UserID  string
	Content string
}

type NoteOutput struct {
	ID        string
	UserID    string
	Content   domain.NoteContent
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListNotesInput struct {
	UserID     string
	Pagination listopts.PaginationInput
	Sort       listopts.SortInput
}

type ListNotesOutput listopts.Paginated[NoteOutput]

type GetUserNoteByIDInput struct {
	UserID string
	ID     string
}

type UpdateNoteInput struct {
	UserID  string
	ID      string
	Content string
}

type DeleteNoteInput struct {
	UserID string
	ID     string
}

type ArchiveNoteInput struct {
	UserID string
	ID     string
}

type RestoreNoteInput struct {
	UserID string
	ID     string
}

func NoteApplicationToOutput(note domain.Note) NoteOutput {
	return NoteOutput{
		ID:        note.ID.Hex(),
		UserID:    note.UserID.Hex(),
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
