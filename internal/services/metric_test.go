package services

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Моки
type MockUpdateMetricStrategy struct {
	mock.Mock
}

func (m *MockUpdateMetricStrategy) Update(metric *domain.Metrics) *domain.Metrics {
	args := m.Called(metric)
	return args.Get(0).(*domain.Metrics)
}

type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) *domain.Metrics {
	args := m.Called(key)
	return args.Get(0).(*domain.Metrics)
}

type MockKeyEncoder struct {
	mock.Mock
}

func (m *MockKeyEncoder) Encode(id, mtype string) string {
	args := m.Called(id, mtype)
	return args.String(0)
}

type MockRanger struct {
	mock.Mock
}

func (m *MockRanger) Range(callback func(key string, value *domain.Metrics) bool) {
	m.Called(callback)
}

func TestUpdateMetricService_Update_GaugeMetric(t *testing.T) {
	mockGaugeStrategy := new(MockUpdateMetricStrategy)
	mockCounterStrategy := new(MockUpdateMetricStrategy)

	// Используем пустой логгер для теста
	mockLogger := zap.NewNop()

	updateService := NewUpdateMetricService(mockCounterStrategy, mockGaugeStrategy, mockLogger)

	value := float64(5)
	metric := &domain.Metrics{ID: "1", MType: domain.Gauge, Value: &value}

	// Ожидаем вызов метода Update на стратегии для Gauge
	mockGaugeStrategy.On("Update", metric).Return(metric)

	result := updateService.UpdateMetric(metric)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockGaugeStrategy.AssertExpectations(t)
}

func TestUpdateMetricService_Update_CounterMetric(t *testing.T) {
	mockGaugeStrategy := new(MockUpdateMetricStrategy)
	mockCounterStrategy := new(MockUpdateMetricStrategy)

	// Используем пустой логгер для теста
	mockLogger := zap.NewNop()

	updateService := NewUpdateMetricService(mockCounterStrategy, mockGaugeStrategy, mockLogger)

	delta := int64(5)
	metric := &domain.Metrics{ID: "1", MType: domain.Counter, Delta: &delta}

	// Ожидаем вызов метода Update на стратегии для Counter
	mockCounterStrategy.On("Update", metric).Return(metric)

	result := updateService.UpdateMetric(metric)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockCounterStrategy.AssertExpectations(t)
}

func TestUpdateMetricService_Update_UnknownMetricType(t *testing.T) {
	mockGaugeStrategy := new(MockUpdateMetricStrategy)
	mockCounterStrategy := new(MockUpdateMetricStrategy)

	// Используем пустой логгер для теста
	mockLogger := zap.NewNop()

	updateService := NewUpdateMetricService(mockCounterStrategy, mockGaugeStrategy, mockLogger)

	// Метрика с неизвестным типом
	metric := &domain.Metrics{ID: "1", MType: "UnknownType"} // "UnknownType" - это строка, что может быть невалидным типом

	// Вызываем UpdateMetric
	result := updateService.UpdateMetric(metric)

	// Проверяем, что результат будет nil
	assert.Nil(t, result)
}

func TestGetMetricService_Get(t *testing.T) {
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	// Используем пустой логгер для теста
	mockLogger := zap.NewNop()

	getMetricService := NewGetMetricService(mockGetter, mockEncoder, mockLogger)

	metric := &domain.Metrics{ID: "1", MType: domain.Gauge, Value: nil}

	mockEncoder.On("Encode", "1", "gauge").Return("gauge_1")
	mockGetter.On("Get", "gauge_1").Return(metric)

	result := getMetricService.GetMetric("1", domain.Gauge)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
}

func TestGetAllMetricsService_GetAll(t *testing.T) {
	mockRanger := new(MockRanger)

	// Используем пустой логгер для теста
	mockLogger := zap.NewNop()

	getAllService := NewGetAllMetricsService(mockRanger, mockLogger)

	metric1 := &domain.Metrics{ID: "1", MType: domain.Gauge, Value: nil}
	metric2 := &domain.Metrics{ID: "2", MType: domain.Counter, Value: nil}

	mockRanger.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(0).(func(string, *domain.Metrics) bool)
		callback("gauge_1", metric1)
		callback("counter_2", metric2)
	}).Return()

	result := getAllService.GetAllMetrics()

	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Contains(t, result, metric1)
	assert.Contains(t, result, metric2)
	mockRanger.AssertExpectations(t)
}
