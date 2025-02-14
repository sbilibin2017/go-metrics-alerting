package responders

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponder_Respond(t *testing.T) {
	// Инициализируем новый Gin-экземпляр для теста
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Создаем группу маршрутов для тестирования
	r.GET("/test-success", func(c *gin.Context) {
		responder := &SuccessResponder{C: c}
		// Симулируем успешный ответ
		responder.Respond("Success", http.StatusOK)
	})

	t.Run("Should return correct success message with 200 status", func(t *testing.T) {
		// Отправляем тестовый запрос к маршруту
		w := performRequest(r, "GET", "/test-success")

		// Проверяем статус код
		assert.Equal(t, http.StatusOK, w.Code)

		// Проверяем содержимое ответа
		assert.Equal(t, "Success", w.Body.String())
	})

	t.Run("Should return custom success message with custom status code", func(t *testing.T) {
		// Добавим обработчик для кастомного сообщения и кода статуса
		r.GET("/test-custom-success", func(c *gin.Context) {
			responder := &SuccessResponder{C: c}
			// Симулируем успешный ответ с другим статусом
			responder.Respond("Custom Success Message", http.StatusCreated)
		})

		// Отправляем запрос к маршруту
		w := performRequest(r, "GET", "/test-custom-success")

		// Проверяем статус код
		assert.Equal(t, http.StatusCreated, w.Code)

		// Проверяем содержимое ответа
		assert.Equal(t, "Custom Success Message", w.Body.String())
	})
}
