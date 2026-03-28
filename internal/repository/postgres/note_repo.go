package postgres

import (
	"api/internal/domain"
	"api/internal/repository/postgres/db"
	"context"
)

type noteRepository struct {
	queries *db.Queries
}

func NewNoteRepository(queries *db.Queries) domain.NoteRepository {
	return &noteRepository{queries: queries}
}

func toDomainNote(row db.Note) *domain.Note {
	return &domain.Note{
		UserID:  row.UserID,
		Title:   row.Title,
		Content: row.Content,
		BaseModel: domain.BaseModel{
			ID:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}
}

func (r *noteRepository) CreateNote(ctx context.Context, userID int64, title, content string) (*domain.Note, error) {
	row, err := r.queries.CreateNote(ctx, db.CreateNoteParams{
		UserID:  userID,
		Title:   title,
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	return toDomainNote(row), nil
}

func (r *noteRepository) DeleteNote(ctx context.Context, id int64) error {
	return r.queries.DeleteNote(ctx, id)
}

func (r *noteRepository) GetNoteByID(ctx context.Context, id int64) (*domain.Note, error) {
	row, err := r.queries.GetNoteByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainNote(row), nil
}

func (r *noteRepository) ListNotes(ctx context.Context) ([]*domain.Note, error) {
	rows, err := r.queries.ListNotes(ctx)
	if err != nil {
		return nil, err
	}
	var notes []*domain.Note
	for _, row := range rows {
		notes = append(notes, toDomainNote(row))
	}
	return notes, nil
}

func (r *noteRepository) ListUserNotes(ctx context.Context, userID int64) ([]*domain.Note, error) {
	rows, err := r.queries.ListUserNotes(ctx, userID)
	if err != nil {
		return nil, err
	}
	var notes []*domain.Note
	for _, row := range rows {
		notes = append(notes, toDomainNote(row))
	}
	return notes, nil
}

func (r *noteRepository) UpdateNote(ctx context.Context, id int64, title, content string) (*domain.Note, error) {
	row, err := r.queries.UpdateNote(ctx, db.UpdateNoteParams{
		ID:      id,
		Title:   title,
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	return toDomainNote(row), nil
}
