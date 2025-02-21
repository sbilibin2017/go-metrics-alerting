package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware - middleware для логирования HTTP-запросов и ответов в Gin
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("uri", c.Request.RequestURI),
			zap.Int("status", statusCode),
			zap.Int("response_size", responseSize),
			zap.Duration("duration", duration),
		)
	}
}
