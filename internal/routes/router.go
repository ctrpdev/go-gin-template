package routes

import (
	"time"

	"api/internal/handler/http"
	"api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter inicializa el router global, aplicando middlewares generales
// y orquestando los diferentes m�dulos de rutas de la aplicaci�n.
func SetupRouter(userHandler *http.UserHandler, noteHandler *http.NoteHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	// Usamos gin.New() en lugar de Default() para no inyectar el logger viejo
	r := gin.New()

	// Paso 5: Configuraci�n estricta de CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://tudominio.com"}, // Or�genes permitidos (tu frontend)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // ESTO ES CR�TICO para recibir las HTTP-Only Cookies
		MaxAge:           12 * time.Hour,
	}))

	// Inyectamos nuestro custom logger estructurado y el Recovery est�ndar de Gin (para panics)
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// Health Check general (buena práctica en APIs modernas para Kubernetes/Docker)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Agrupar rutas (opcionalmente bajo un prefijo como /api/v1)
	api := r.Group("/")

	// Delegar la configuración de rutas de usuario a su propio archivo
	SetupUserRoutes(api, userHandler, authMiddleware)

	// Rutas de Notas
	SetupNoteRoutes(api, noteHandler, authMiddleware)

	return r
}
