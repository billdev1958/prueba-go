package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger es el middleware adaptado para Gin de uno que era originalmente para la libreria estandar
func GinLogger(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start).Seconds()
		status := c.Writer.Status()

		logEntry := slog.Group("request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.Float64("latency_seconds", latency),
			slog.Int("status", status),
			slog.String("user_agent", c.Request.UserAgent()),
		)

		if status >= 500 {
			l.Error("request failed", logEntry)
		} else {
			l.Info("request completed", logEntry)
		}
	}
}
