package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterMetricRoutes(t *testing.T) {
	// Создаем новый роутер Gin
	router := gin.Default()

	// Заглушки для обработчиков
	dummyHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}

	// Регистрируем маршруты с заглушками
	RegisterMetricRoutes(
		router,
		dummyHandler,
		dummyHandler,
		dummyHandler,
		dummyHandler,
		dummyHandler,
	)

	// Создаем тестовые запросы и проверяем ответы
	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/update/"},
		{"POST", "/update/gauge/test/123"},
		{"POST", "/value/"},
		{"GET", "/value/gauge/test"},
		{"GET", "/"},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "ok")
		})
	}
}
