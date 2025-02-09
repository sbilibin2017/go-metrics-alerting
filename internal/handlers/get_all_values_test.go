package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок-сервис для GetAllValuesService
type MockGetAllValuesService struct {
	mock.Mock
}

func (m *MockGetAllValuesService) GetAllMetricValues(ctx context.Context) []*types.MetricResponse {
	args := m.Called(ctx)
	if result, ok := args.Get(0).([]*types.MetricResponse); ok {
		return result
	}
	return nil
}

func setupGetAllTestRouter(svc GetAllValuesService) *gin.Engine {
	r := gin.Default()
	RegisterGetAllMetricValuesHandler(r, svc)
	return r
}

func TestGetAllMetricValuesHandler_Success(t *testing.T) {
	// Arrange
	mockService := new(MockGetAllValuesService)
	router := setupGetAllTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetAllMetricValues", mock.Anything).Return([]*types.MetricResponse{
		{UpdateMetricValueRequest: types.UpdateMetricValueRequest{Type: "gauge", Name: "cpu", Value: "100"}},
		{UpdateMetricValueRequest: types.UpdateMetricValueRequest{Type: "gauge", Name: "memory", Value: "2048"}},
	})

	// Создаем запрос
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<h1>Metrics List</h1>")
	assert.Contains(t, w.Body.String(), "<li>cpu: 100</li>")
	assert.Contains(t, w.Body.String(), "<li>memory: 2048</li>")

	// Проверяем заголовки
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("Date"))

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetAllMetricValuesHandler_NoMetrics(t *testing.T) {
	// Arrange
	mockService := new(MockGetAllValuesService)
	router := setupGetAllTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetAllMetricValues", mock.Anything).Return(nil)

	// Создаем запрос
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error": "No metrics found"}`, w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}
