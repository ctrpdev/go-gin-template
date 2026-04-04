package routes

import (
	"api/internal/handler/http"
	"api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupNoteRoutes registra todas las rutas correspondientes al dominio de Notas
func SetupNoteRoutes(rg *gin.RouterGroup, noteHandler *http.NoteHandler, authMiddleware *middleware.AuthMiddleware) {
	// Rutas Protegidas de Notas
	// Todas requieren autenticación
	notes := rg.Group("/notes")
	notes.Use(authMiddleware.RequireAuth())
	{
		notes.POST("/", noteHandler.CreateNote)
		notes.GET("/", noteHandler.ListUserNotes)
		notes.GET("/all", noteHandler.ListNotes) // Puede requerir rol especial después
		notes.GET("/:id", noteHandler.GetNoteByID)
		notes.PUT("/:id", noteHandler.UpdateNote)
		notes.DELETE("/:id", noteHandler.DeleteNote)
	}
}
