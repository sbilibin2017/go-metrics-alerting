package services

import (
	"go-metrics-alerting/internal/metric/handlers"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetricStorage struct {
	mock.Mock
}

func (m *MockMetricStorage) Save(metricType string, metricName string, value string) error {
	args := m.Called(metricType, metricName, value)
	return args.Error(0)
}

func (m *MockMetricStorage) Get(metricType string, metricName string) (string, error) {
	args := m.Called(metricType, metricName)
	return args.String(0), args.Error(1)
}

func TestUpdateMetricValueService_UpdateMetricValue_Success(t *testing.T) {
	mockStorage := new(MockMetricStorage)
	service := NewUpdateMetricValueService(mockStorage)

	// Мокируем получение текущего значения метрики
	mockStorage.On("Get", string(GaugeType), "cpu").Return("10", nil)

	// Мокируем сохранение обновленного значения
	mockStorage.On("Save", string(GaugeType), "cpu", "20").Return(nil)

	req := &handlers.UpdateMetricValueRequest{
		Type:  string(GaugeType),
		Name:  "cpu",
		Value: "20",
	}

	// Тестируем успешное обновление
	err := service.UpdateMetricValue(req)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestUpdateMetricValueService_UpdateMetricValue_GaugeError(t *testing.T) {
	mockStorage := new(MockMetricStorage)
	service := NewUpdateMetricValueService(mockStorage)

	// Мокируем получение текущего значения метрики
	mockStorage.On("Get", string(GaugeType), "cpu").Return("10", nil)

	req := &handlers.UpdateMetricValueRequest{
		Type:  string(GaugeType),
		Name:  "cpu",
		Value: "invalid", // Неправильное значение
	}

	// Тестируем ошибку при обновлении значения метрики
	err := service.UpdateMetricValue(req)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, ErrInvalidGaugeValue.Error(), apiErr.Message)
}

func TestUpdateMetricValueService_UpdateMetricValue_CounterSuccess(t *testing.T) {
	mockStorage := new(MockMetricStorage)
	service := NewUpdateMetricValueService(mockStorage)

	// Мокируем получение текущего значения метрики
	mockStorage.On("Get", string(CounterType), "requests").Return("10", nil)

	// Мокируем сохранение обновленного значения
	mockStorage.On("Save", string(CounterType), "requests", "20").Return(nil)

	req := &handlers.UpdateMetricValueRequest{
		Type:  string(CounterType),
		Name:  "requests",
		Value: "10", // Прибавляем 10
	}

	// Тестируем успешное обновление счетчика
	err := service.UpdateMetricValue(req)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestUpdateMetricValueService_UpdateMetricValue_CounterError(t *testing.T) {
	mockStorage := new(MockMetricStorage)
	service := NewUpdateMetricValueService(mockStorage)

	// Мокируем получение текущего значения метрики
	mockStorage.On("Get", string(CounterType), "requests").Return("10", nil)

	req := &handlers.UpdateMetricValueRequest{
		Type:  string(CounterType),
		Name:  "requests",
		Value: "invalid", // Неправильное значение
	}

	// Тестируем ошибку при обновлении значения метрики
	err := service.UpdateMetricValue(req)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, ErrInvalidCounterValue.Error(), apiErr.Message)
}

func TestUpdateMetricValueService_UpdateMetricValue_MetricMismatch(t *testing.T) {
	mockStorage := new(MockMetricStorage)
	service := NewUpdateMetricValueService(mockStorage)

	// Мокируем получение текущего значения метрики
	mockStorage.On("Get", string(CounterType), "requests").Return("10", nil)

	req := &handlers.UpdateMetricValueRequest{
		Type:  string(CounterType),
		Name:  "requests",
		Value: "-15", // Приведет к переполнению
	}

	// Тестируем ошибку при переполнении счетчика
	err := service.UpdateMetricValue(req)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, ErrMetricMismatch.Error(), apiErr.Message)
}

func TestUpdateMetricValueService_UpdateMetricValue_UnsupportedMetricType(t *testing.T) {
	mockStorage := new(MockMetricStorage)
	service := NewUpdateMetricValueService(mockStorage)

	// Мокируем вызов Get для любого типа метрики
	mockStorage.On("Get", mock.Anything, mock.Anything).Return("", nil)

	req := &handlers.UpdateMetricValueRequest{
		Type:  "unknown", // Неподдерживаемый тип метрики
		Name:  "requests",
		Value: "10",
	}

	// Тестируем ошибку для неподдерживаемого типа метрики
	err := service.UpdateMetricValue(req)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, ErrUnsupportedMetricType.Error(), apiErr.Message)
}
