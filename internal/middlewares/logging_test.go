package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLoggingMiddleware(t *testing.T) {
	// Используем zaptest.NewLogger для перехвата логов
	logger := zaptest.NewLogger(t)

	// Инициализируем тестовую обработку запроса
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// Создаем middleware
	loggingMiddleware := LoggingMiddleware(logger)

	// Применяем middleware к обработчику
	handlerWithMiddleware := loggingMiddleware(handler)

	// Создаем новый запрос
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос с middleware
	handlerWithMiddleware.ServeHTTP(w, req)

	// Проверяем, что ответ был правильным
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, world!", w.Body.String())

	// Проверяем, что в логе появились записи о запросе и ответе
	// Проверка запроса
	logger.Check(zap.InfoLevel, "Request").Write(zap.String("method", "GET"), zap.String("uri", "/hello"))

	// Проверка ответа
	logger.Check(zap.InfoLevel, "Response").Write(
		zap.Int("status_code", http.StatusOK),
		zap.Int64("response_size", int64(len("Hello, world!"))),
	)
}

func TestLoggingMiddleware_DurationLogged(t *testing.T) {
	// Используем zaptest.NewLogger для перехвата логов
	logger := zaptest.NewLogger(t)

	// Инициализируем тестовую обработку запроса
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Искусственная задержка
		w.Write([]byte("Hello, world!"))
	})

	// Создаем middleware
	loggingMiddleware := LoggingMiddleware(logger)

	// Применяем middleware к обработчику
	handlerWithMiddleware := loggingMiddleware(handler)

	// Создаем новый запрос
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос с middleware
	handlerWithMiddleware.ServeHTTP(w, req)

	// Проверяем, что ответ был правильным
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, world!", w.Body.String())

	// Проверяем, что в логе была запись о длительности
	logger.Check(zap.InfoLevel, "Response").Write(
		zap.Int("status_code", http.StatusOK),
		zap.Int64("response_size", int64(len("Hello, world!"))),
		zap.Duration("duration", 100*time.Millisecond),
	)
}
