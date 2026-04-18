package app

import (
	"time"

	"momento/internal/notes/domain"
	"momento/pkg/listopts"
)

type NoteInput struct {
	UserID  string
	Title   string
	Content string
}

type NoteOutput struct {
	ID        string
	UserID    string
	Title     domain.NoteTitle
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
	Title   *string
	Content *string
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
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
