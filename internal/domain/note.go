package domain

import (
	"context"
)

type Note struct {
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	BaseModel
}

type NoteRepository interface {
	CreateNote(ctx context.Context, userID int64, title, content string) (*Note, error)
	DeleteNote(ctx context.Context, id int64) error
	GetNoteByID(ctx context.Context, id int64) (*Note, error)
	ListNotes(ctx context.Context) ([]*Note, error)
	ListUserNotes(ctx context.Context, userID int64) ([]*Note, error)
	UpdateNote(ctx context.Context, id int64, title, content string) (*Note, error)
}

type NoteService interface {
	CreateNote(ctx context.Context, userID int64, title, content string) (*Note, error)
	DeleteNote(ctx context.Context, id int64) error
	GetNoteByID(ctx context.Context, id int64) (*Note, error)
	ListNotes(ctx context.Context) ([]*Note, error)
	ListUserNotes(ctx context.Context, userID int64) ([]*Note, error)
	UpdateNote(ctx context.Context, id int64, title, content string) (*Note, error)
}
