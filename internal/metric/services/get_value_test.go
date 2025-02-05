package services

import (
	"errors"
	"go-metrics-alerting/internal/metric/handlers"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для интерфейса MetricStorageGetter
type MockMetricStorageGetter struct {
	mock.Mock
}

func (m *MockMetricStorageGetter) Get(metricType string, metricName string) (string, error) {
	args := m.Called(metricType, metricName)
	return args.String(0), args.Error(1)
}

func TestGetMetricValueService_GetMetricValue_Success(t *testing.T) {
	mockStorage := new(MockMetricStorageGetter)
	service := NewGetMetricValueService(mockStorage)

	// Мокируем успешный ответ для метрики
	mockStorage.On("Get", "gauge", "cpu").Return("99", nil)

	req := &handlers.GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}

	// Тестируем успешный случай получения метрики
	value, err := service.GetMetricValue(req)

	assert.NoError(t, err)
	assert.Equal(t, "99", value)
	mockStorage.AssertExpectations(t)
}

func TestGetMetricValueService_GetMetricValue_MetricNotFound(t *testing.T) {
	mockStorage := new(MockMetricStorageGetter)
	service := NewGetMetricValueService(mockStorage)

	// Мокируем ошибку, что метрика не найдена
	mockStorage.On("Get", "gauge", "cpu").Return("", errors.New("metric not found"))

	req := &handlers.GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}

	// Тестируем ошибку, когда метрика не найдена
	value, err := service.GetMetricValue(req)

	assert.Equal(t, "", value)
	assert.Error(t, err)

	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
	assert.Equal(t, "metric not found", apiErr.Message)

	mockStorage.AssertExpectations(t)
}
