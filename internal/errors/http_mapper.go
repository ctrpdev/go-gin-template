package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ResponseError es nuestro formato centralizado para enviar JSON al cliente
type ResponseError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// MapDomainError intercepta los errores y mapea "Dominio" -> "Status Code" de Gin
func MapDomainError(c *gin.Context, err error) {
	// 1. Manejar Errores de Validación (Tags binding: required, min, email) del JSON
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]string, len(ve))
		for i, fe := range ve {
			out[i] = "Failed field " + fe.Field() + " on " + fe.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": out})
		return
	}

	// 2. Manejar Errores Conocidos de Dominio
	switch {
	case errors.Is(err, ErrUserNotFound):
		c.JSON(http.StatusNotFound, ResponseError{Error: err.Error()})
	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ResponseError{Error: err.Error()})
	case errors.Is(err, ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, ResponseError{Error: err.Error()})
	case errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrTokenRevoked):
		c.JSON(http.StatusUnauthorized, ResponseError{Error: err.Error()})

	// 3. Cualquier error sin clasificar -> 500
	default:
		c.JSON(http.StatusInternalServerError, ResponseError{Error: "Internal Server Error", Message: err.Error()})
	}
}
