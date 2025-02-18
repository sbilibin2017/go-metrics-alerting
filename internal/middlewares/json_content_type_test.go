package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJSONContentTypeMiddleware_ValidContentType(t *testing.T) {
	// Устанавливаем режим тестирования
	gin.SetMode(gin.TestMode)

	// Создаем новый маршрутизатор
	r := gin.Default()

	// Регистрируем мидлвару и тестовый маршрут
	r.Use(JSONContentTypeMiddleware())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Тестируем запрос с правильным Content-Type
	reqBody := []byte(`{"key":"value"}`)
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(recorder, req)

	// Проверяем, что статус ответа — 200 (OK)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Проверяем, что тело ответа содержит ожидаемое сообщение
	assert.Contains(t, recorder.Body.String(), `"message":"success"`)
}

func TestJSONContentTypeMiddleware_InvalidContentType(t *testing.T) {
	// Устанавливаем режим тестирования
	gin.SetMode(gin.TestMode)

	// Создаем новый маршрутизатор
	r := gin.Default()

	// Регистрируем мидлвару и тестовый маршрут
	r.Use(JSONContentTypeMiddleware())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Тестируем запрос с неправильным Content-Type
	reqBody := []byte(`{"key":"value"}`)
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	recorder := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(recorder, req)

	// Проверяем, что статус ответа — 415 (Unsupported Media Type)
	assert.Equal(t, http.StatusUnsupportedMediaType, recorder.Code)

	// Проверяем, что тело ответа содержит ошибку "Unsupported Media Type"
	assert.Contains(t, recorder.Body.String(), `"error":"Unsupported Media Type"`)
}
