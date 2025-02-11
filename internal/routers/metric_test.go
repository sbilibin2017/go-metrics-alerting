package routers

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

// Мок-сервис для MetricService
type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) UpdateMetric(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockMetricService) GetMetric(ctx context.Context, req *types.GetMetricValueRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockMetricService) ListMetrics(ctx context.Context) []*types.MetricResponse {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MetricResponse)
}

// Тест для роута /update/:type/:name/:value
func TestRegisterMetricHandlers_UpdateMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Создаем мок-сервис
	mockSvc := new(MockMetricService)

	// Регистрация роута
	RegisterMetricHandlers(r, mockSvc)

	// Ожидаем, что при вызове UpdateMetric будет возвращен nil (ошибки не будет)
	mockSvc.On("UpdateMetric", mock.Anything, &types.UpdateMetricValueRequest{
		Type:  "cpu",
		Name:  "usage",
		Value: "75",
	}).Return(nil)

	// Отправляем POST запрос
	req, _ := http.NewRequest("POST", "/update/cpu/usage/75", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что метод был вызван
	mockSvc.AssertExpectations(t)
}

// Тест для роута /value/:type/:name
func TestRegisterMetricHandlers_GetMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Создаем мок-сервис
	mockSvc := new(MockMetricService)

	// Регистрация роута
	RegisterMetricHandlers(r, mockSvc)

	// Ожидаем, что метод GetMetric вернет строку "75" без ошибки
	mockSvc.On("GetMetric", mock.Anything, &types.GetMetricValueRequest{
		Type: "cpu",
		Name: "usage",
	}).Return("75", nil)

	// Отправляем GET запрос
	req, _ := http.NewRequest("GET", "/value/cpu/usage", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем тело ответа
	assert.Equal(t, "75", w.Body.String())

	// Проверяем, что метод был вызван
	mockSvc.AssertExpectations(t)
}

// Тест для роута /
func TestRegisterMetricHandlers_ListMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Создаем мок-сервис
	mockSvc := new(MockMetricService)

	// Регистрация роута
	RegisterMetricHandlers(r, mockSvc)

	// Ожидаем, что метод ListMetrics вернет срез с одним элементом
	mockSvc.On("ListMetrics", mock.Anything).Return([]*types.MetricResponse{
		{
			Name:  "usage",
			Value: "75",
		},
	})

	// Отправляем GET запрос
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, w.Code)

	// Получаем фактический ответ сервера
	actualResponse := w.Body.String()

	// Проверяем, что ответ содержит "usage: 75"
	assert.Contains(t, actualResponse, "usage: 75")

	// Проверяем, что метод был вызван
	mockSvc.AssertExpectations(t)
}
