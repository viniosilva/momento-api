package application

import (
	"pinnado/internal/notes/domain"
	shareddto "pinnado/internal/shared/application/dto"
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
	Pagination shareddto.PaginationInput
	Sort       shareddto.SortInput
}

type ListNotesOutput shareddto.Paginated[NoteOutput]

func NoteApplicationToOutput(note domain.Note) NoteOutput {
	return NoteOutput{
		ID:        note.ID.Hex(),
		UserID:    note.UserID.Hex(),
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
