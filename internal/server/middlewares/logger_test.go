package middlewares

import (
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger - Мок для интерфейса Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

// TestLoggerMiddleware_CheckDuration проверяет, что метод Info логирует продолжительность запроса
func TestLoggerMiddleware_CheckDuration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Создаем мок логгера
	mockLogger := new(MockLogger)

	// Устанавливаем ожидания для вызова метода Info
	mockLogger.On("Info", "Request processed", mock.MatchedBy(func(fields []zap.Field) bool {
		for _, field := range fields {
			if field.Key == "duration" {
				// Проверяем, что продолжительность больше нуля
				if field.Integer > 0 {
					return true
				}
			}
		}
		return false
	})).Return(nil)

	// Регистрируем middleware
	r.Use(LoggerMiddleware(mockLogger))

	// Регистрируем маршрут
	r.GET("/test", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.String(http.StatusOK, "OK")
	})

	// Создаем запрос
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	// Обрабатываем запрос
	r.ServeHTTP(recorder, req)

	// Проверяем, что метод Info был вызван с duration
	mockLogger.AssertExpectations(t)
}
