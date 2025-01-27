package responders

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponder реализует ErrorResponderInterface
type ErrorResponder struct{}

// NewErrorResponder создает новый ErrorResponder
func NewErrorResponder() *ErrorResponder {
	return &ErrorResponder{}
}

// Respond теперь выполняет обработку ошибки, передавая контекст
func (er *ErrorResponder) Respond(c *gin.Context, statusCode int, message string) {
	// Устанавливаем заголовок Content-Type
	c.Header("Content-Type", "text/plain; charset=utf-8")
	// Отправляем текстовое сообщение ошибки
	c.String(statusCode, message)
}
