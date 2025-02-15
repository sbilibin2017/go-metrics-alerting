package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware логирует запросы и ответы.
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Считываем данные запроса
		start := time.Now()

		// Запоминаем URI и метод
		uri := c.FullPath()
		method := c.Request.Method

		// Выполняем обработку запроса
		c.Next()

		// Вычисляем время выполнения запроса
		duration := time.Since(start)

		// Логируем данные
		logger.Info("Request processed",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.Duration("duration", duration),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("content_length", c.Writer.Size()),
		)
	}
}
