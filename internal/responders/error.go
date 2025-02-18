package responders

import (
	"go-metrics-alerting/internal/logger" // Импортируем логгер

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RespondWithError отправляет JSON-ответ с ошибкой и логирует ее.
func RespondWithError(c *gin.Context, statusCode int, err error) {
	// Логируем ошибку, которая произошла
	logger.Logger.Error("Error occurred in request",
		zap.Int("status_code", statusCode),
		zap.String("error_message", err.Error()))

	// Отправляем JSON-ответ с ошибкой
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
