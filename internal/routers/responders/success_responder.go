package responders

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SuccessResponder реализует SuccessResponderInterface
type SuccessResponder struct{}

// NewSuccessResponder создает новый экземпляр SuccessResponder
func NewSuccessResponder() *SuccessResponder {
	return &SuccessResponder{}
}

// Отправка текстового ответа с установкой правильных заголовков
func (sr *SuccessResponder) Respond(c *gin.Context, statusCode int, message string) {
	// Устанавливаем статус ответа
	c.Status(statusCode)

	// Устанавливаем заголовок Content-Type
	c.Header("Content-Type", "text/plain; charset=utf-8")

	// Если сообщение пустое, устанавливаем Content-Length = 0
	if message == "" {
		c.Header("Content-Length", "0")
	} else {
		// Устанавливаем Content-Length на основе длины сообщения
		c.Header("Content-Length", strconv.Itoa(len(message)))
		// Отправляем текстовое сообщение
		c.String(statusCode, message)
	}

	// Устанавливаем заголовок Date
	c.Header("Date", time.Now().UTC().Format(http.TimeFormat))
}
