package application

import (
	"context"
	"errors"
	"fmt"

	"pinnado/internal/notes/domain"
	"pinnado/pkg/listopts"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteService struct {
	noteRepository NoteRepository
}

func NewNoteService(noteRepository NoteRepository) *NoteService {
	return &NoteService{
		noteRepository: noteRepository,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, input NoteInput) (NoteOutput, error) {
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("invalid user ID: %w", err)
	}

	title, err := domain.NewNoteTitle(input.Title)
	if err != nil {
		return NoteOutput{}, err
	}

	content, err := domain.NewNoteContent(input.Content)
	if err != nil {
		return NoteOutput{}, err
	}

	note := domain.NewNote(userID, title, content)

	if err := s.noteRepository.Create(ctx, note); err != nil {
		return NoteOutput{}, fmt.Errorf("s.noteRepository.Create: %w", err)
	}

	return NoteApplicationToOutput(note), nil
}

func (s *NoteService) ListNotes(ctx context.Context, input ListNotesInput) (ListNotesOutput, error) {
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return ListNotesOutput{}, fmt.Errorf("invalid user ID: %w", err)
	}

	params := listopts.ListParams{
		Pagination: input.Pagination,
		Sort:       input.Sort,
	}
	paginatedNotes, err := s.noteRepository.ListByUserID(ctx, userID, params)
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

func (s *NoteService) GetUserNoteByID(ctx context.Context, input GetUserNoteByIDInput) (NoteOutput, error) {
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("invalid ID: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("invalid user ID: %w", err)
	}

	note, err := s.noteRepository.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return NoteOutput{}, domain.ErrNoteNotFound
		}

		return NoteOutput{}, fmt.Errorf("s.noteRepository.GetByIDAndUserID: %w", err)
	}

	noteOutput := NoteApplicationToOutput(note)
	return noteOutput, nil
}

func (s *NoteService) UpdateNote(ctx context.Context, input UpdateNoteInput) (NoteOutput, error) {
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("invalid ID: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("invalid user ID: %w", err)
	}

	title, err := domain.NewNoteTitle(input.Title)
	if err != nil {
		return NoteOutput{}, err
	}

	content, err := domain.NewNoteContent(input.Content)
	if err != nil {
		return NoteOutput{}, err
	}

	note, err := s.noteRepository.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("s.noteRepository.GetByIDAndUserID: %w", err)
	}
	note.SetTitle(title)
	note.SetContent(content)

	if err := s.noteRepository.Update(ctx, note); err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return NoteOutput{}, domain.ErrNoteNotFound
		}

		return NoteOutput{}, fmt.Errorf("s.noteRepository.Update: %w", err)
	}

	return NoteApplicationToOutput(note), nil
}

func (s *NoteService) DeleteNote(ctx context.Context, input DeleteNoteInput) error {
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := s.noteRepository.DeleteByIDAndUserID(ctx, id, userID); err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return domain.ErrNoteNotFound
		}

		return fmt.Errorf("s.noteRepository.DeleteByIDAndUserID: %w", err)
	}

	return nil
}

func (s *NoteService) ArchiveNote(ctx context.Context, input ArchiveNoteInput) error {
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.noteRepository.ArchiveByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return domain.ErrNoteNotFound
		}

		return fmt.Errorf("s.noteRepository.ArchiveByIDAndUserID: %w", err)
	}

	return nil
}

func (s *NoteService) RestoreNote(ctx context.Context, input RestoreNoteInput) error {
	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.noteRepository.RestoreByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return domain.ErrNoteNotFound
		}

		return fmt.Errorf("s.noteRepository.RestoreByIDAndUserID: %w", err)
	}

	return nil
}
