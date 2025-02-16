package middlewares

import (
	"time"

	"go-metrics-alerting/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware logs requests and responses using the global logger.
func LoggerMiddleware() gin.HandlerFunc {
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

		// Log the request details using the global logger
		logger.Log.Info("Request processed",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Duration("duration", duration),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("content_length", c.Writer.Size()),
		)
	}
}
