package responders

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HTMLResponder реализует HTMLResponderInterface
type HTMLResponder struct{}

// NewHTMLResponder создает новый экземпляр HTMLResponder
func NewHTMLResponder() *HTMLResponder {
	return &HTMLResponder{}
}

// Отправка HTML-страницы с установкой правильных заголовков
func (hr *HTMLResponder) Respond(c *gin.Context, statusCode int, message string) {
	// Устанавливаем статус и тип содержимого
	c.Data(statusCode, "text/html; charset=utf-8", []byte(message))

	// Устанавливаем заголовок Date
	c.Header("Date", time.Now().UTC().Format(http.TimeFormat))
}
