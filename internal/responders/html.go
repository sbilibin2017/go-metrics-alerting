package responders

import (
	"net/http"
	"text/template"

	"go-metrics-alerting/internal/logger" // Импортируем логгер

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RespondWithHTML рендерит HTML-шаблон и отправляет его клиенту, добавлено логирование.
func RespondWithHTML(c *gin.Context, statusCode int, tmplString string, data interface{}) {
	// Логируем начало обработки запроса на рендеринг HTML
	logger.Logger.Info("Rendering HTML response",
		zap.Int("status_code", statusCode),
		zap.String("template", tmplString),
		zap.Any("data", data))

	// Парсим HTML-шаблон
	tmpl, err := template.New("response").Parse(tmplString)
	if err != nil {
		// Логируем ошибку, если не удалось распарсить шаблон
		logger.Logger.Error("Failed to parse HTML template", zap.Error(err))
		RespondWithError(c, http.StatusInternalServerError, err)
		return
	}

	// Отправляем заголовок с типом контента
	c.Header("Content-Type", "text/html")
	c.Writer.WriteHeader(statusCode)

	// Выполняем шаблон
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		// Логируем ошибку, если не удалось выполнить шаблон
		logger.Logger.Error("Failed to execute HTML template", zap.Error(err))
		RespondWithError(c, http.StatusInternalServerError, err)
		return
	}

	// Логируем успешный рендеринг ответа
	logger.Logger.Info("HTML response rendered successfully")
}
