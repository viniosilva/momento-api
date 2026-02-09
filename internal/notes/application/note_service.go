package application

import (
	"context"
	"fmt"

	"pinnado/internal/notes/domain"

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

	return NoteOutput{
		ID:        note.ID.Hex(),
		UserID:    note.UserID.Hex(),
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}
