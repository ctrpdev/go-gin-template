package routes

import (
	"api/internal/handler/http"
	"api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes registra todas las rutas correspondientes al dominio de Usuarios
func SetupUserRoutes(rg *gin.RouterGroup, userHandler *http.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// Rutas Públicas de Usuario
	rg.POST("/register", userHandler.Register)
	rg.POST("/login", userHandler.Login)
	rg.POST("/refresh", userHandler.Refresh)
	rg.POST("/verify/:id", userHandler.Verify)

	// Rutas Protegidas de Usuario
	protected := rg.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/logout", userHandler.Logout)
		protected.GET("/me", userHandler.GetMe)
	}
}
