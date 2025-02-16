package services

import (
	"errors"
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetter реализует интерфейс Getter для тестирования
type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

// MockSetter реализует интерфейс Setter для тестирования
type MockSetter struct {
	mock.Mock
}

func (m *MockSetter) Set(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

// MockRanger реализует интерфейс Ranger для тестирования
type MockRanger struct {
	mock.Mock
}

func (m *MockRanger) Range(callback func(key string, value string) bool) {
	m.Called(callback)
}

// --- Тесты для UpdateMetricService ---

func TestUpdateMetric_Counter_Success(t *testing.T) {
	mockGetter := new(MockGetter)
	mockSetter := new(MockSetter)
	service := NewUpdateMetricService(mockGetter, mockSetter)

	mockGetter.On("Get", "counter_metric").Return("10", nil)
	mockSetter.On("Set", "counter_metric", "20").Return(nil)

	delta := int64(10)
	req := &types.MetricsRequest{
		ID:    "counter_metric",
		MType: types.Counter,
		Delta: &delta,
	}

	resp, err := service.UpdateMetric(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(20), *resp.Delta)
	mockGetter.AssertExpectations(t)
	mockSetter.AssertExpectations(t)
}

func TestUpdateMetric_Counter_NotFound(t *testing.T) {
	mockGetter := new(MockGetter)
	mockSetter := new(MockSetter)
	service := NewUpdateMetricService(mockGetter, mockSetter)

	mockGetter.On("Get", "counter_metric").Return("", errors.New("not found"))
	mockSetter.On("Set", "counter_metric", "10").Return(nil)

	delta := int64(10)
	req := &types.MetricsRequest{
		ID:    "counter_metric",
		MType: types.Counter,
		Delta: &delta,
	}

	resp, err := service.UpdateMetric(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(10), *resp.Delta)
	mockGetter.AssertExpectations(t)
	mockSetter.AssertExpectations(t)
}

func TestUpdateMetric_Gauge_Success(t *testing.T) {
	mockGetter := new(MockGetter)
	mockSetter := new(MockSetter)
	service := NewUpdateMetricService(mockGetter, mockSetter)

	// Ожидаем вызова метода Get, чтобы предотвратить ошибку
	mockGetter.On("Get", "gauge_metric").Return("0", nil)
	mockSetter.On("Set", "gauge_metric", "15.5").Return(nil)

	value := 15.5
	req := &types.MetricsRequest{
		ID:    "gauge_metric",
		MType: types.Gauge,
		Value: &value,
	}

	resp, err := service.UpdateMetric(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 15.5, *resp.Value)
	mockGetter.AssertExpectations(t)
	mockSetter.AssertExpectations(t)
}

// --- Тесты для GetMetricService ---

func TestGetMetric_Success(t *testing.T) {
	mockGetter := new(MockGetter)
	service := NewGetMetricService(mockGetter)

	mockGetter.On("Get", "metric_1").Return("42", nil)

	req := &types.MetricsRequest{ID: "metric_1"}
	value, err := service.GetMetric(req)

	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "42", *value)
	mockGetter.AssertExpectations(t)
}

func TestGetMetric_NotFound(t *testing.T) {
	mockGetter := new(MockGetter)
	service := NewGetMetricService(mockGetter)

	mockGetter.On("Get", "metric_1").Return("", errors.New("not found"))

	req := &types.MetricsRequest{ID: "metric_1"}
	value, err := service.GetMetric(req)

	assert.Error(t, err)
	assert.Nil(t, value)
	assert.Equal(t, ErrMetricIsNotFound, err)
	mockGetter.AssertExpectations(t)
}

// --- Тесты для ListMetricsService ---

func TestListMetrics_Success(t *testing.T) {
	mockRanger := new(MockRanger)
	service := NewListMetricsService(mockRanger)

	mockRanger.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(0).(func(string, string) bool)
		callback("metric_1", "42")
		callback("metric_2", "3.14")
	}).Return()

	metrics := service.ListMetrics()

	assert.Len(t, metrics, 2)
	assert.Equal(t, "metric_1", metrics[0].ID)
	assert.Equal(t, "42", metrics[0].Value)
	assert.Equal(t, "metric_2", metrics[1].ID)
	assert.Equal(t, "3.14", metrics[1].Value)
	mockRanger.AssertExpectations(t)
}

func TestUpdateMetric_SetterError(t *testing.T) {
	mockGetter := new(MockGetter)
	mockSetter := new(MockSetter)
	service := NewUpdateMetricService(mockGetter, mockSetter)

	// Настройка для Get, чтобы возвращать значение "0" для метрики
	mockGetter.On("Get", "counter_metric").Return("0", nil)

	// Настройка для Set, чтобы он возвращал ошибку
	mockSetter.On("Set", "counter_metric", "10").Return(errors.New("DB error"))

	// Запрос на обновление метрики
	delta := int64(10)
	req := &types.MetricsRequest{
		ID:    "counter_metric",
		MType: types.Counter,
		Delta: &delta,
	}

	// Выполнение метода
	resp, err := service.UpdateMetric(req)

	// Проверка, что произошла ошибка
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrMetricIsNotUpdated, err)

	// Проверка, что все ожидания на моки выполнены
	mockGetter.AssertExpectations(t)
	mockSetter.AssertExpectations(t)
}

func TestUpdateMetric_Gauge_SetterError(t *testing.T) {
	mockGetter := new(MockGetter)
	mockSetter := new(MockSetter)
	service := NewUpdateMetricService(mockGetter, mockSetter)

	// Настройка для Get, чтобы возвращать значение "0" для метрики
	mockGetter.On("Get", "gauge_metric").Return("0", nil)

	// Настройка для Set, чтобы он возвращал ошибку
	mockSetter.On("Set", "gauge_metric", "15.5").Return(errors.New("DB error"))

	// Запрос на обновление метрики
	value := 15.5
	req := &types.MetricsRequest{
		ID:    "gauge_metric",
		MType: types.Gauge,
		Value: &value,
	}

	// Выполнение метода
	resp, err := service.UpdateMetric(req)

	// Проверка, что произошла ошибка
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrMetricIsNotUpdated, err)

	// Проверка, что все ожидания на моки выполнены
	mockGetter.AssertExpectations(t)
	mockSetter.AssertExpectations(t)
}
