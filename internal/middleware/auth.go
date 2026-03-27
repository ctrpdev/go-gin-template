package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	secretKey   []byte
	sessionRepo domain.SessionRepository
}

func NewAuthMiddleware(secretKey []byte, sessionRepo domain.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey:   secretKey,
		sessionRepo: sessionRepo,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid format"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.secretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Check if blacklisted
		jti, ok := claims["jti"].(string)
		if ok {
			blacklisted, err := m.sessionRepo.IsAccessTokenBlacklisted(c.Request.Context(), jti)
			if err != nil || blacklisted {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
				return
			}
			c.Set("jti", jti)
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Subject missing in token"})
			return
		}

		userID, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid subject in token"})
			return
		}

		// Inject role
		role, _ := claims["role"].(string)
		if role == "" {
			role = "user" // Default fallback
		}

		// Inject to context
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}

// RequireRole es un middleware que verifica que el usuario tenga un rol especifico
func (m *AuthMiddleware) RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		for _, role := range requiredRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
	}
}
