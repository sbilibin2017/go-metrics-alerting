package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// Middleware для установки заголовков, включая Content-Type и Content-Length
func SetHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем новый буфер для записи ответа
		var buf bytes.Buffer
		// Создаем новый ResponseWriter, который будет писать в буфер
		writer := &responseSetHeadersWriter{ResponseWriter: w, buffer: &buf}

		// Определяем Content-Type, например, на основе URL или других параметров
		contentType := "application/json" // По умолчанию
		if r.URL.Path == "/text" {        // Пример: для пути "/text" используем text/plain
			contentType = "text/plain"
		} else if r.URL.Path == "/html" { // Пример: для пути "/html" используем text/html
			contentType = "text/html"
		}

		// Устанавливаем заголовок Content-Type
		writer.Header().Set("Content-Type", contentType)

		// Устанавливаем заголовок Date в формате RFC1123
		writer.Header().Set("Date", time.Now().UTC().Format(time.RFC1123))

		// Передаем запрос в следующий обработчик
		next.ServeHTTP(writer, r)

		// Теперь у нас есть полный ответ в буфере, можно установить Content-Length
		contentLength := len(buf.Bytes())
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))

		// Теперь записываем ответ из буфера в оригинальный ResponseWriter
		w.Write(buf.Bytes())
	})
}

// responseWriter - расширяем стандартный http.ResponseWriter, чтобы захватывать ответ в буфер
type responseSetHeadersWriter struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func (rw *responseSetHeadersWriter) Write(p []byte) (n int, err error) {
	return rw.buffer.Write(p)
}
