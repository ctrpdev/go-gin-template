package http

import (
	"net/http"
	"strconv"

	"api/internal/domain"
	"api/internal/errors"
	"api/internal/handler/http/dto"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	noteService domain.NoteService
}

func NewNoteHandler(noteService domain.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.noteService.CreateNote(c.Request.Context(), userID, req.Title, req.Content)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": note})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	err = h.noteService.DeleteNote(c.Request.Context(), id)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *NoteHandler) GetNoteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	note, err := h.noteService.GetNoteByID(c.Request.Context(), id)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}

func (h *NoteHandler) ListNotes(c *gin.Context) {
	notes, err := h.noteService.ListNotes(c.Request.Context())
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func (h *NoteHandler) ListUserNotes(c *gin.Context) {
	userID := c.GetInt64("userID")

	notes, err := h.noteService.ListUserNotes(c.Request.Context(), userID)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req dto.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.noteService.UpdateNote(c.Request.Context(), id, req.Title, req.Content)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}
