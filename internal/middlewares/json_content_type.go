package middlewares

import (
	"net/http"

	"go-metrics-alerting/internal/logger" // Импортируем логгер

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	JsonContentType         string = "application/json"
	ErrUnsupportedMediaType string = "Unsupported Media Type"
)

// JSONContentTypeMiddleware проверяет, что Content-Type = application/json для POST-запросов и логирует информацию о запросах
func JSONContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Логируем информацию о входящем запросе
		logger.Logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()))

		if c.Request.Method == http.MethodPost {
			// Логируем проверку Content-Type
			if c.ContentType() != JsonContentType {
				// Логируем ошибку с неверным Content-Type
				logger.Logger.Warn("Unsupported Media Type", zap.String("url", c.Request.URL.String()))
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": ErrUnsupportedMediaType})
				c.Abort()
				return
			}
		}

		// Передаем управление следующему обработчику
		c.Next()

		// Логируем статус ответа
		logger.Logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Int("status_code", c.Writer.Status()))
	}
}
