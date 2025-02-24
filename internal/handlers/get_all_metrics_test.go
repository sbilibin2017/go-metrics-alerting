package handlers

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок-сервис для тестирования
type MockGetAllMetricsService struct {
	mock.Mock
}

func (m *MockGetAllMetricsService) GetAllMetrics() []*domain.Metrics {
	args := m.Called()
	return args.Get(0).([]*domain.Metrics)
}

func TestGetAllMetricsHandler_Success(t *testing.T) {
	// Подготовим тестовые данные
	mockService := new(MockGetAllMetricsService)

	// Создадим метрики для теста
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			MType: domain.Counter,
			Delta: new(int64),
		},
		{
			ID:    "metric2",
			MType: domain.Gauge,
			Value: new(float64),
		},
	}

	// Мокаем вызов GetAllMetrics
	mockService.On("GetAllMetrics").Return(metrics)

	// Создаем HTTP-запрос для тестирования
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	// Регистрируем обработчик с мок-сервисом
	handler := GetAllMetricsHandler(mockService)
	handler.ServeHTTP(w, req)

	// Проверяем правильность ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что метод GetAllMetrics был вызван
	mockService.AssertExpectations(t)
}

func TestGetAllMetricsHandler_NoMetrics(t *testing.T) {
	// Подготовим пустую коллекцию метрик для теста
	mockService := new(MockGetAllMetricsService)

	// Мокаем вызов GetAllMetrics, чтобы он вернул пустой список
	mockService.On("GetAllMetrics").Return([]*domain.Metrics{})

	// Создаем HTTP-запрос для тестирования
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	// Регистрируем обработчик с мок-сервисом
	handler := GetAllMetricsHandler(mockService)
	handler.ServeHTTP(w, req)

	// Проверяем правильность ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что шаблон был вызван и что в ответе нет метрик
	expectedResponse := `<div><p>ID: </p><p>Value: </p></div>`
	assert.NotContains(t, w.Body.String(), expectedResponse)

	// Проверяем, что метод GetAllMetrics был вызван
	mockService.AssertExpectations(t)
}
