package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger defines the methods that a logger must implement.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

// LoggerMiddleware logs requests and responses using the injected logger.
func LoggerMiddleware(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time
		start := time.Now()

		// Capture the URI and method
		uri := c.FullPath()
		method := c.Request.Method

		// Process the request
		c.Next()

		// Calculate the request duration
		duration := time.Since(start)

		// Log the request details
		log.Info("Request processed",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Duration("duration", duration),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("content_length", c.Writer.Size()),
		)
	}
}
