package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger - интерфейс для логгера
type Logger interface {
	Info(msg string, fields ...zap.Field)
}

func LoggerMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		uri := c.FullPath()
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)

		logger.Info("Request processed",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Duration("duration", duration),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("content_length", c.Writer.Size()),
		)
	}
}
