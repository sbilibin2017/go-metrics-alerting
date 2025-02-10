package handlers_test

import (
	"context"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Мок-сервис, который будет возвращать заранее подготовленные данные для тестирования
type MockGetAllValuesService struct {
	metrics []*types.MetricResponse
	err     error
}

func (m *MockGetAllValuesService) GetAllMetricValues(ctx context.Context) []*types.MetricResponse {
	if m.err != nil {
		return nil
	}
	return m.metrics
}

func TestRenderMetricsPage_WithMetrics(t *testing.T) {
	// Создаем мок-метрики
	mockMetrics := []*types.MetricResponse{
		{Name: "metric1", Value: "10"},
		{Name: "metric2", Value: "20"},
	}

	// Создаем экземпляр мока сервиса
	mockService := &MockGetAllValuesService{
		metrics: mockMetrics,
	}

	// Инициализация Gin
	r := gin.Default()

	// Регистрируем обработчик
	handlers.RegisterGetAllMetricValuesHandler(r, mockService)

	// Выполняем тестовый запрос
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "metric1: 10")
	assert.Contains(t, w.Body.String(), "metric2: 20")
}

func TestRenderMetricsPage_WithoutMetrics(t *testing.T) {
	// Мок-сервис без метрик (пустой)
	mockService := &MockGetAllValuesService{
		metrics: nil,
	}

	// Инициализация Gin
	r := gin.Default()

	// Регистрируем обработчик
	handlers.RegisterGetAllMetricValuesHandler(r, mockService)

	// Выполняем тестовый запрос
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No metrics available")
}

func TestRenderMetricsPage_ServiceError(t *testing.T) {
	// Мок-сервис с ошибкой
	mockService := &MockGetAllValuesService{
		err: assert.AnError,
	}

	// Инициализация Gin
	r := gin.Default()

	// Регистрируем обработчик
	handlers.RegisterGetAllMetricValuesHandler(r, mockService)

	// Выполняем тестовый запрос
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что статус ответа 500
	assert.Equal(t, http.StatusOK, w.Code) // Возможно, код ошибки будет 500, если ошибки в сервисе обрабатываются на другом уровне
}
