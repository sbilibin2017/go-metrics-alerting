package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestLoggingMiddleware(t *testing.T) {
	// Создаем мок-логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Отложенная синхронизация логгера

	// Создаем тестовый обработчик, который будет использоваться в middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Применяем middleware
	loggingHandler := LoggingMiddleware(logger, testHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Вызов middleware
	loggingHandler.ServeHTTP(w, req)

	// Проверяем код статуса ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем размер ответа
	assert.Equal(t, len("Hello, World!"), w.Body.Len())

}

func TestLoggingMiddlewareLogsRequestDetails(t *testing.T) {
	// Создаем буфер для захвата логов
	var buf zaptest.Buffer

	// Создаем mock логгер, который будет записывать логи в буфер
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	))

	// Создаем тестовый обработчик
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test"))
	})

	// Применяем middleware
	loggingHandler := LoggingMiddleware(logger, testHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Вызов middleware
	loggingHandler.ServeHTTP(w, req)

	// Проверка, что запрос прошел успешно
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка, что логи содержат все ключевые данные
	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP Response")
	assert.Contains(t, logOutput, "method")
	assert.Contains(t, logOutput, "uri")
	assert.Contains(t, logOutput, "status")
	assert.Contains(t, logOutput, "response_size")
	assert.Contains(t, logOutput, "duration")
}
