package responders

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithHTML_Success(t *testing.T) {
	// Настроим тестовый режим для Gin
	gin.SetMode(gin.TestMode)

	// Создаем новый ResponseRecorder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	// Тестовые данные
	tmplString := "<html><body><h1>{{.Title}}</h1></body></html>"
	data := map[string]interface{}{
		"Title": "Test Page",
	}

	// Вызов RespondWithHTML
	RespondWithHTML(c, http.StatusOK, tmplString, data)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Проверяем, что в теле ответа есть заголовок
	expectedHTML := "<html><body><h1>Test Page</h1></body></html>"
	assert.Equal(t, expectedHTML, recorder.Body.String())

	// Если необходимо, можно добавить проверку логов
	// Например, можно проверить, что в логах есть запись о рендеринге шаблона.
}

func TestRespondWithHTML_TemplateParseError(t *testing.T) {
	// Настроим тестовый режим для Gin
	gin.SetMode(gin.TestMode)

	// Создаем новый ResponseRecorder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	// Тест с неправильным шаблоном
	tmplString := "<html><body><h1>{{.Title}</h1></body></html>" // Ошибка в шаблоне

	// Вызов RespondWithHTML с ошибкой
	RespondWithHTML(c, http.StatusInternalServerError, tmplString, nil)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Проверяем, что ошибка была залогирована
	// Здесь можно использовать проверку на содержание в логах
}

func TestRespondWithHTML_TemplateExecutionError(t *testing.T) {
	// Настроим тестовый режим для Gin
	gin.SetMode(gin.TestMode)

	// Создаем новый ResponseRecorder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	// Тест с правильно парсированным шаблоном, но ошибкой при выполнении
	tmplString := "<html><body><h1>{{.Title}}</h1></body></html>"

	// Вызов RespondWithHTML с некорректными данными
	RespondWithHTML(c, http.StatusInternalServerError, tmplString, make(chan int)) // Ошибка выполнения шаблона

	// Проверяем статус ответа
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

}
