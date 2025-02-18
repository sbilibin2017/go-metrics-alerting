package responders

import (
	"go-metrics-alerting/internal/logger" // Импортируем логгер

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// respondWithSuccess отправляет JSON-ответ с успешными данными и логирует информацию о запросе.
func RespondWithSuccess(c *gin.Context, statusCode int, payload interface{}) {
	// Логируем отправку успешного ответа
	logger.Logger.Info("Responding with success",
		zap.Int("status_code", statusCode),
		zap.Any("payload", payload))

	// Отправляем JSON-ответ с успешными данными
	c.JSON(statusCode, payload)
}
