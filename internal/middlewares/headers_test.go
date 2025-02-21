package middlewares

import (
	"fmt"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Функция для создания тестового роута с middleware
func createTestRouter() *gin.Engine {
	r := gin.Default()
	r.Use(SetHeadersMiddleware())

	// Пример маршрута, который возвращает JSON
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// Пример маршрута, который возвращает текст
	r.GET("/text", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Plain Text!")
	})

	return r
}

func TestSetHeadersMiddleware_JSONResponse(t *testing.T) {
	// Создаем тестовый роутер
	r := createTestRouter()

	// Выполняем запрос и получаем ответ через httptest.NewRecorder
	w := performRequest(r, "GET", "/json")

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем заголовки
	assert.Contains(t, w.Header(), "Date")
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf("%d", len(`{"message":"Hello, World!"}`)), w.Header().Get("Content-Length"))
}

func TestSetHeadersMiddleware_TextResponse(t *testing.T) {
	// Создаем тестовый роутер
	r := createTestRouter()

	// Выполняем запрос и получаем ответ через httptest.NewRecorder
	w := performRequest(r, "GET", "/text")

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем заголовки
	assert.Contains(t, w.Header(), "Date")
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf("%d", len("Hello, Plain Text!")), w.Header().Get("Content-Length"))
}

func TestSetHeadersMiddleware_ExistingContentTypeHeader(t *testing.T) {
	// Создаем тестовый роутер
	r := createTestRouter()

	// Выполняем запрос с заголовком Content-Type
	req, _ := http.NewRequest("GET", "/json", nil)
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	w := httptest.NewRecorder()

	// Обрабатываем запрос с middleware
	r.ServeHTTP(w, req)

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем заголовки
	assert.Contains(t, w.Header(), "Date")
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf("%d", len(`{"message":"Hello, World!"}`)), w.Header().Get("Content-Length"))
}

// performRequest выполняет HTTP-запрос и возвращает httptest.ResponseRecorder
func performRequest(r *gin.Engine, method, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
