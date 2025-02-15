package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware logs requests and responses.
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
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
		logger.Info("Request processed",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Duration("duration", duration),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("content_length", c.Writer.Size()),
		)
	}
}
