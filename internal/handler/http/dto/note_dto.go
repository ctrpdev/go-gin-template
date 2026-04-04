package dto

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required,max=100"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title" binding:"omitempty,max=100"`
	Content string `json:"content" binding:"omitempty"`
}

type NoteResponse struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListNotesResponse struct {
	Notes []NoteResponse `json:"notes"`
}
