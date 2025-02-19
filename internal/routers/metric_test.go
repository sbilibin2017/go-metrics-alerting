package routers

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для сервиса обновления метрик
type MockUpdateMetricsService struct {
	mock.Mock
}

func (m *MockUpdateMetricsService) Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*types.UpdateMetricsResponse), args.Error(1)
}

// Мок для сервиса получения значения метрики
type MockGetMetricValueService struct {
	mock.Mock
}

func (m *MockGetMetricValueService) GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*types.GetMetricValueResponse), args.Error(1)
}

// Мок для сервиса получения всех метрик
type MockGetAllMetricValuesService struct {
	mock.Mock
}

func (m *MockGetAllMetricValuesService) GetAllMetricValues() []*types.GetMetricValueResponse {
	args := m.Called()
	return args.Get(0).([]*types.GetMetricValueResponse)
}

func TestRegisterMetricRoutes(t *testing.T) {
	// Создаем маршрутизатор
	r := gin.Default()

	// Моки для сервисов
	mockUpdateService := new(MockUpdateMetricsService)
	mockGetMetricValueService := new(MockGetMetricValueService)
	mockGetAllMetricValuesService := new(MockGetAllMetricValuesService)

	// Регистрируем маршруты
	RegisterMetricRoutes(r, mockUpdateService, mockGetMetricValueService, mockGetAllMetricValuesService)

	// Тестируем POST /update/
	t.Run("POST /update/", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/update/", nil)
		recorder := httptest.NewRecorder()

		// Выполняем запрос
		r.ServeHTTP(recorder, req)

		// Проверяем, что статус ответа — 404 (метод не найден)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
	})

	// Тестируем POST /value/
	t.Run("POST /value/", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/value/", nil)
		recorder := httptest.NewRecorder()

		// Выполняем запрос
		r.ServeHTTP(recorder, req)

		// Проверяем, что статус ответа — 404 (метод не найден)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
	})

	// Настроим мок для GetAllMetricValues
	mockGetAllMetricValuesService.On("GetAllMetricValues").Return([]*types.GetMetricValueResponse{
		{ID: "1", Value: "10"},
		{ID: "2", Value: "20"},
	})

	// Тестируем GET /
	t.Run("GET /", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// Выполняем запрос
		r.ServeHTTP(recorder, req)

		// Проверяем, что статус ответа — 200 (OK)
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Проверяем, что ответ содержит данные метрик
		assert.Contains(t, recorder.Body.String(), "Metrics List")
		assert.Contains(t, recorder.Body.String(), "<li>1: 10</li>")
		assert.Contains(t, recorder.Body.String(), "<li>2: 20</li>")

		// Проверяем, что мок был вызван
		mockGetAllMetricValuesService.AssertCalled(t, "GetAllMetricValues")
	})
}
