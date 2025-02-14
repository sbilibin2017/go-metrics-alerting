package responders

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/test-success", handleSuccess)
	r.GET("/test-parse-error", handleParseError)
	r.GET("/test-execution-error", handleExecutionError)

	return r
}

// Тест 1: Успешный рендеринг HTML-шаблона
func TestHTMLHandler_Success(t *testing.T) {
	r := setupRouter()
	w := performRequest(r, "GET", "/test-success")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<h1>Hello, world!</h1>")
}

// Тест 2: Ошибка при парсинге шаблона
func TestHTMLHandler_ParseError(t *testing.T) {
	r := setupRouter()
	w := performRequest(r, "GET", "/test-parse-error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Template parsing failed")
}

// Обработчик для успешного рендеринга HTML
func handleSuccess(c *gin.Context) {
	handler := &HTMLHandler{C: c}
	templateContent := "<html><body><h1>{{.Title}}</h1></body></html>"
	data := map[string]string{"Title": "Hello, world!"}
	handler.Respond(templateContent, data)
}

// Обработчик для ошибки парсинга шаблона
func handleParseError(c *gin.Context) {
	handler := &HTMLHandler{C: c}
	// Ошибка в синтаксисе шаблона
	templateContent := "<html><body><h1>{{.Title}</h1></body></html>" // Пропущена закрывающая скобка
	data := map[string]string{"Title": "Hello, world!"}
	handler.Respond(templateContent, data)
}

// Обработчик для ошибки выполнения шаблона
func handleExecutionError(c *gin.Context) {
	handler := &HTMLHandler{C: c}
	templateContent := "<html><body><h1>{{.Title}}</h1></body></html>"
	// Передаем неподдерживаемый тип данных, который вызовет ошибку
	handler.Respond(templateContent, make(chan int))
}
