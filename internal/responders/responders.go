package responders

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Enum для типов респондера
type ResponderType string

const (
	JSONResponder   ResponderType = "json"
	StringResponder ResponderType = "string"
	HTMLResponder   ResponderType = "html"
)

// Функция-ответчик для различных форматов
func Respond(c *gin.Context, responderType ResponderType, statusCode int, payload interface{}) {
	setHeaders(c, responderType)

	switch responderType {
	case JSONResponder:
		c.JSON(statusCode, payload)
	case StringResponder:
		c.String(statusCode, payload.(string))
	case HTMLResponder:
		renderHTML(c, statusCode, payload)
	}
}

// Рендеринг HTML-страницы
func renderHTML(c *gin.Context, statusCode int, payload interface{}) {
	tmplString, ok := payload.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid HTML template"})
		return
	}

	tmpl, err := template.New("response").Parse(tmplString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(statusCode)
	err = tmpl.Execute(c.Writer, c.Keys["metrics"]) // Передаем метрики в шаблон
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Установка заголовков
func setHeaders(c *gin.Context, responderType ResponderType) {
	switch responderType {
	case JSONResponder:
		c.Header("Content-Type", "application/json; charset=utf-8")
	case StringResponder:
		c.Header("Content-Type", "text/plain; charset=utf-8")
	case HTMLResponder:
		c.Header("Content-Type", "text/html; charset=utf-8")
	}
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
}
