package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware - Middleware для логирования запросов и ответов
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := &LoggingResponseWriter{ResponseWriter: w}
			// Логирование запроса
			logger.Info("Request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Time("start_time", start),
			)
			next.ServeHTTP(lrw, r)

			// Логирование ответа
			duration := time.Since(start)
			logger.Info("Response",
				zap.Int("status_code", lrw.statusCode),
				zap.Int64("response_size", lrw.bodySize),
				zap.Duration("duration", duration),
			)
		})
	}
}

// LoggingResponseWriter - Обертка для ResponseWriter, чтобы захватывать статус код и размер ответа
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bodySize   int64
}

// WriteHeader переопределяет метод WriteHeader для захвата статус кода
func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	if lrw.statusCode == 0 { // Устанавливаем статус код только если еще не был установлен
		lrw.statusCode = statusCode
	}
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Write переопределяет метод Write для захвата размера тела ответа
func (lrw *LoggingResponseWriter) Write(p []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(p)
	lrw.bodySize += int64(size)
	return size, err
}
