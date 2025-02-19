package responders

import (
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
)

// Общая функция для установки заголовков
func setHeaders(c *gin.Context, contentType string) {
	c.Header("Content-Type", contentType)
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
}

// Вспомогательные функции

func SendErrorJSON(c *gin.Context, statusCode int, message string) {
	setHeaders(c, "application/json; charset=utf-8")
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

func SendSuccessJSON(c *gin.Context, statusCode int, response interface{}) {
	setHeaders(c, "application/json; charset=utf-8")
	c.JSON(statusCode, response)
}

func SendSuccessText(c *gin.Context, statusCode int, message string) {
	setHeaders(c, "text/plain; charset=utf-8")
	c.String(statusCode, message)
}

// Функция для отправки HTML
func SendSuccessHTML(c *gin.Context, statusCode int, templateString string, data interface{}) {
	setHeaders(c, "text/html; charset=utf-8")

	// Парсим шаблон
	tmpl, err := template.New("response").Parse(templateString)
	if err != nil {
		c.String(500, "Error parsing template: "+err.Error())
		return
	}

	// Отправляем HTML-страницу
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		c.String(500, "Error rendering template: "+err.Error())
	}
}
