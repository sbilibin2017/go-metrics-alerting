package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware - middleware для логирования HTTP-запросов и ответов в Chi
func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем ResponseWriter для захвата статуса и размера ответа
		writer := &responseLoggingWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Передаем запрос в следующий обработчик
		next.ServeHTTP(writer, r)

		duration := time.Since(start)
		statusCode := writer.statusCode
		responseSize := writer.size

		// Логируем информацию о статусе, размере и времени ответа
		logger.Info("HTTP Response",
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path), // Используем r.URL.Path для корректного пути
			zap.Int("status", statusCode),
			zap.Int("response_size", responseSize),
			zap.Duration("duration", duration),
		)
	})
}

// responseLoggingWriter - расширяем стандартный http.ResponseWriter, чтобы захватывать статус и размер ответа
type responseLoggingWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseLoggingWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ResponseWriter.Write(p)
	rw.size += n
	return n, err
}

func (rw *responseLoggingWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
