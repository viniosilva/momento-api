package app

import (
	"context"
	"errors"
	"fmt"

	"momento/internal/notes/domain"
	"momento/pkg/listopts"
	"momento/pkg/tools"
)

type noteService struct {
	noteRepository NoteRepository
}

func NewNoteService(noteRepository NoteRepository) *noteService {
	return &noteService{
		noteRepository: noteRepository,
	}
}

func (s *noteService) CreateNote(ctx context.Context, input NoteInput) (NoteOutput, error) {
	title, err := domain.NewNoteTitle(input.Title)
	if err != nil {
		return NoteOutput{}, err
	}

	content, err := domain.NewNoteContent(input.Content)
	if err != nil {
		return NoteOutput{}, err
	}

	note := domain.NewNote(input.UserID, title, content)

	if err := s.noteRepository.Create(ctx, note); err != nil {
		return NoteOutput{}, fmt.Errorf("s.noteRepository.Create: %w", err)
	}

	return NoteApplicationToOutput(note), nil
}

func (s *noteService) ListNotes(ctx context.Context, input ListNotesInput) (ListNotesOutput, error) {
	params := listopts.ListParams{
		Pagination: input.Pagination,
		Sort:       input.Sort,
	}
	paginatedNotes, err := s.noteRepository.ListByUserID(ctx, input.UserID, params)
	if err != nil {
		return ListNotesOutput{}, fmt.Errorf("s.noteRepository.ListByUserID: %w", err)
	}

	noteOutputs := make([]NoteOutput, len(paginatedNotes.Data))
	for i, note := range paginatedNotes.Data {
		noteOutputs[i] = NoteApplicationToOutput(note)
	}

	return ListNotesOutput{
		Data:       noteOutputs,
		Pagination: paginatedNotes.Pagination,
	}, nil
}

func (s *noteService) GetUserNoteByID(ctx context.Context, input GetUserNoteByIDInput) (NoteOutput, error) {
	note, err := s.noteRepository.GetByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return NoteOutput{}, domain.ErrNoteNotFound
		}

		return NoteOutput{}, fmt.Errorf("s.noteRepository.GetByIDAndUserID: %w", err)
	}

	return NoteApplicationToOutput(note), nil
}

func (s *noteService) UpdateNote(ctx context.Context, input UpdateNoteInput) (NoteOutput, error) {
	var newTitle *domain.NoteTitle
	if input.Title != nil {
		title, err := domain.NewNoteTitle(*input.Title)
		if err != nil {
			return NoteOutput{}, err
		}
		newTitle = &title
	}

	var newContent *domain.NoteContent
	if input.Content != nil {
		content, err := domain.NewNoteContent(*input.Content)
		if err != nil {
			return NoteOutput{}, err
		}
		newContent = &content
	}

	note, err := s.noteRepository.GetByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("s.noteRepository.GetByIDAndUserID: %w", err)
	}

	note.Update(
		tools.ValueOrDefault(newTitle, note.Title),
		tools.ValueOrDefault(newContent, note.Content),
	)

	if err := s.noteRepository.Update(ctx, note); err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return NoteOutput{}, domain.ErrNoteNotFound
		}

		return NoteOutput{}, fmt.Errorf("s.noteRepository.Update: %w", err)
	}

	return NoteApplicationToOutput(note), nil
}

func (s *noteService) DeleteNote(ctx context.Context, input DeleteNoteInput) error {
	if err := s.noteRepository.DeleteByIDAndUserID(ctx, input.ID, input.UserID); err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return domain.ErrNoteNotFound
		}

		return fmt.Errorf("s.noteRepository.DeleteByIDAndUserID: %w", err)
	}

	return nil
}

func (s *noteService) ArchiveNote(ctx context.Context, input ArchiveNoteInput) error {
	err := s.noteRepository.ArchiveByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return domain.ErrNoteNotFound
		}

		return fmt.Errorf("s.noteRepository.ArchiveByIDAndUserID: %w", err)
	}

	return nil
}

func (s *noteService) RestoreNote(ctx context.Context, input RestoreNoteInput) error {
	err := s.noteRepository.RestoreByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return domain.ErrNoteNotFound
		}

		return fmt.Errorf("s.noteRepository.RestoreByIDAndUserID: %w", err)
	}

	return nil
}
