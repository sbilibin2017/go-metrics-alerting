package responders

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithSuccess(t *testing.T) {
	// Настроим тестовый режим для Gin
	gin.SetMode(gin.TestMode)

	// Создаем новый ResponseRecorder
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	// Тестовые данные
	payload := map[string]interface{}{
		"message": "success",
	}

	// Вызов RespondWithSuccess
	RespondWithSuccess(c, http.StatusOK, payload)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Проверяем тело ответа
	expectedJSON := `{"message":"success"}`
	assert.JSONEq(t, expectedJSON, recorder.Body.String())

	// Здесь вы можете добавить проверку для логирования, если нужно.
	// Примерно так:
	// assert.Contains(t, yourLogBuffer, "Responding with success")
}
