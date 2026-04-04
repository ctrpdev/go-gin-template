package service

import (
	"api/internal/domain"
	"context"
)

type noteService struct {
	noteRepo domain.NoteRepository
}

func NewNoteService(noteRepo domain.NoteRepository) domain.NoteService {
	return &noteService{
		noteRepo: noteRepo,
	}
}

func (s *noteService) CreateNote(ctx context.Context, userID int64, title, content string) (*domain.Note, error) {
	return s.noteRepo.CreateNote(ctx, userID, title, content)
}

func (s *noteService) DeleteNote(ctx context.Context, id int64) error {
	return s.noteRepo.DeleteNote(ctx, id)
}

func (s *noteService) GetNoteByID(ctx context.Context, id int64) (*domain.Note, error) {
	return s.noteRepo.GetNoteByID(ctx, id)
}

func (s *noteService) ListNotes(ctx context.Context) ([]*domain.Note, error) {
	return s.noteRepo.ListNotes(ctx)
}

func (s *noteService) ListUserNotes(ctx context.Context, userID int64) ([]*domain.Note, error) {
	return s.noteRepo.ListUserNotes(ctx, userID)
}

func (s *noteService) UpdateNote(ctx context.Context, id int64, title, content string) (*domain.Note, error) {
	return s.noteRepo.UpdateNote(ctx, id, title, content)
}
