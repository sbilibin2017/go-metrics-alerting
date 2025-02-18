package responders

import (
	"errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithError(t *testing.T) {
	// Настроим тестовый режим для Gin
	gin.SetMode(gin.TestMode)

	// Создаем новый ResponseRecorder, который реализует http.ResponseWriter
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	// Создаем ошибку
	err := errors.New("some error occurred")

	// Вызов RespondWithError
	RespondWithError(c, http.StatusBadRequest, err)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Проверяем тело ответа
	expectedJSON := `{"error": "some error occurred"}`
	assert.JSONEq(t, expectedJSON, recorder.Body.String())
}
