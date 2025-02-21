package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware для установки заголовков
func SetHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Создаем новый ResponseWriter для перехвата тела ответа
		writer := &ResponseWriter{ResponseWriter: c.Writer, Body: []byte{}}
		c.Writer = writer

		// Обрабатываем запрос
		c.Next()

		// Получаем тело ответа и устанавливаем заголовки после обработки запроса
		body := string(writer.Body) // Преобразуем тело в строку
		var contentType string      // По умолчанию пусто

		// Определяем тип контента в зависимости от тела ответа
		if c.GetHeader("Content-Type") == "" {
			// Если тело содержит JSON
			if len(body) > 0 && (body[0] == '{' || body[0] == '[') {
				contentType = "application/json; charset=utf-8"
			} else {
				contentType = "text/plain; charset=utf-8"
			}
		} else {
			contentType = c.GetHeader("Content-Type")
		}

		// Устанавливаем общие заголовки
		c.Header("Date", time.Now().UTC().Format(time.RFC1123))
		c.Header("Content-Type", contentType)
		c.Header("Content-Length", fmt.Sprintf("%d", len(body)))
	}
}

// ResponseWriter - структура для захвата тела ответа
type ResponseWriter struct {
	gin.ResponseWriter
	Body []byte
}

// Write - перезаписываем метод Write, чтобы захватывать тело ответа
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...) // Сохраняем данные в body
	return w.ResponseWriter.Write(b)
}
