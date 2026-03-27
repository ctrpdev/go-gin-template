package errors

import "errors"

// Errores de Dominio Comunes
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("email already registered")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenRevoked       = errors.New("token has been revoked")
)
