package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger es un middleware que registra cada petición HTTP usando slog
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Pasar a los siguientes handlers
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Ignorar health check en los logs (opcional, para no spamear)
		if path == "/health" {
			return
		}

		if status >= 400 && status < 500 {
			slog.Warn("HTTP Request",
				slog.String("method", c.Request.Method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("ip", c.ClientIP()),
				slog.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
		} else if status >= 500 {
			slog.Error("HTTP Request",
				slog.String("method", c.Request.Method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("ip", c.ClientIP()),
				slog.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
		} else {
			slog.Info("HTTP Request",
				slog.String("method", c.Request.Method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("ip", c.ClientIP()),
			)
		}
	}
}
