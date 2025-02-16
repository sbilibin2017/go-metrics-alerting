package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware_RealServer(t *testing.T) {
	// Инициализируем реальный zap логгер
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // Ensure that any buffered log entries are flushed

	// Создаем новый Gin рутер
	r := gin.New()

	// Добавляем middleware
	r.Use(LoggerMiddleware())

	// Добавляем маршрут
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Создаем новый HTTP запрос
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	// Результат выполнения запроса
	w := httptest.NewRecorder()

	// Отправляем запрос
	r.ServeHTTP(w, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем тело ответа
	assert.Contains(t, w.Body.String(), `"message":"ok"`)

	// Проверяем логи (это будет вывод в консоль, но можно настроить и проверку через файлы или другие способы)
	// В реальных тестах можно использовать стереотипы или перехватчики для проверки консольных выводов.
}
