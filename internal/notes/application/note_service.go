package application

import (
	"context"
	"fmt"

	"pinnado/internal/notes/domain"
	sharedinfra "pinnado/internal/shared/infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return NoteOutput{}, fmt.Errorf("invalid user ID: %w", err)
	}

	content, err := domain.NewNoteContent(input.Content)
	if err != nil {
		return NoteOutput{}, err
	}

	note := domain.NewNote(userID, content)

	if err := s.noteRepository.Create(ctx, note); err != nil {
		return NoteOutput{}, fmt.Errorf("s.noteRepository.Create: %w", err)
	}

	return NoteApplicationToOutput(note), nil
}

func (s *noteService) ListNotes(ctx context.Context, input ListNotesInput) (ListNotesOutput, error) {
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return ListNotesOutput{}, fmt.Errorf("invalid user ID: %w", err)
	}

	params := sharedinfra.ListParams{
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
