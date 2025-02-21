package services

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(key string, value *domain.Metrics) {
	m.Called(key, value)
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

func TestUpdateCounterMetricService_Update_NewMetric(t *testing.T) {
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	counterService := NewUpdateCounterMetricService(mockSaver, mockGetter, mockEncoder)
	delta := int64(5)
	metric := &domain.Metrics{ID: "1", MType: domain.Counter, Delta: &delta}

	mockEncoder.On("Encode", "1", "counter").Return("counter_1")
	mockGetter.On("Get", "counter_1").Return((*domain.Metrics)(nil)) // Возвращаем nil как типизированное значение
	mockSaver.On("Save", "counter_1", metric).Return()

	result := counterService.Update(metric)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
}

func TestUpdateCounterMetricService_Update_ExistingMetric(t *testing.T) {
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	counterService := NewUpdateCounterMetricService(mockSaver, mockGetter, mockEncoder)
	initialDelta := int64(5)
	existingMetric := &domain.Metrics{ID: "1", MType: domain.Counter, Delta: &initialDelta}
	delta := int64(10)
	metric := &domain.Metrics{ID: "1", MType: domain.Counter, Delta: &delta}

	mockEncoder.On("Encode", "1", "counter").Return("counter_1")
	mockGetter.On("Get", "counter_1").Return(existingMetric)
	mockSaver.On("Save", "counter_1", existingMetric).Return()

	result := counterService.Update(metric)

	assert.NotNil(t, result)
	assert.Equal(t, existingMetric, result)
	assert.Equal(t, int64(15), *existingMetric.Delta) // проверяем, что delta обновился
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
}

func TestUpdateGaugeMetricService_Update_ExistingMetric(t *testing.T) {
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	gaugeService := NewUpdateGaugeMetricService(mockSaver, mockGetter, mockEncoder)

	// Исходная метрика в хранилище
	initialValue := 5.0
	existingMetric := &domain.Metrics{
		ID:    "1",
		MType: domain.Gauge,
		Value: &initialValue,
	}

	// Новая метрика для обновления
	newValue := 10.0
	metric := &domain.Metrics{
		ID:    "1",
		MType: domain.Gauge,
		Value: &newValue,
	}

	// Мокируем поведение
	mockEncoder.On("Encode", "1", "gauge").Return("gauge_1")
	mockGetter.On("Get", "gauge_1").Return(existingMetric)   // Возвращаем уже существующую метрику
	mockSaver.On("Save", "gauge_1", existingMetric).Return() // Сохраняем обновленную метрику

	// Обновляем метрику
	result := gaugeService.Update(metric)

	// Проверяем результат
	assert.NotNil(t, result)
	assert.Equal(t, existingMetric, result)          // Убедитесь, что результат совпадает с обновленной метрикой
	assert.Equal(t, newValue, *existingMetric.Value) // Проверяем, что значение метрики обновилось

	// Проверяем ожидания
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
}

func TestUpdateMetricService_Update_GaugeMetric(t *testing.T) {
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	gaugeService := NewUpdateGaugeMetricService(mockSaver, mockGetter, mockEncoder)
	counterService := NewUpdateCounterMetricService(mockSaver, mockGetter, mockEncoder)
	updateService := NewUpdateMetricService(gaugeService, counterService)

	v := 10.0
	metric := &domain.Metrics{ID: "1", MType: domain.Gauge, Value: &v}

	mockEncoder.On("Encode", "1", "gauge").Return("gauge_1")
	mockGetter.On("Get", "gauge_1").Return((*domain.Metrics)(nil)) // Возвращаем nil как типизированное значение
	mockSaver.On("Save", "gauge_1", metric).Return()

	result := updateService.UpdateMetric(metric)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
}

func TestUpdateMetricService_Update_CounterMetric(t *testing.T) {
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	gaugeService := NewUpdateGaugeMetricService(mockSaver, mockGetter, mockEncoder)
	counterService := NewUpdateCounterMetricService(mockSaver, mockGetter, mockEncoder)
	updateService := NewUpdateMetricService(gaugeService, counterService)

	delta := int64(5)
	metric := &domain.Metrics{ID: "1", MType: domain.Counter, Delta: &delta}

	mockEncoder.On("Encode", "1", "counter").Return("counter_1")
	mockGetter.On("Get", "counter_1").Return((*domain.Metrics)(nil)) // Возвращаем nil как типизированное значение
	mockSaver.On("Save", "counter_1", metric).Return()

	result := updateService.UpdateMetric(metric)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
	mockSaver.AssertExpectations(t)
}

func TestGetMetricService_Get(t *testing.T) {
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	getMetricService := NewGetMetricService(mockGetter, mockEncoder)

	metric := &domain.Metrics{ID: "1", MType: domain.Gauge, Value: nil}

	mockEncoder.On("Encode", "1", "gauge").Return("gauge_1")
	mockGetter.On("Get", "gauge_1").Return(metric)

	result := getMetricService.Get("1", domain.Gauge)

	assert.NotNil(t, result)
	assert.Equal(t, metric, result)
	mockEncoder.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
}

func TestGetAllMetricsService_GetAll(t *testing.T) {
	mockRanger := new(MockRanger)

	getAllService := NewGetAllMetricsService(mockRanger)

	metric1 := &domain.Metrics{ID: "1", MType: domain.Gauge, Value: nil}
	metric2 := &domain.Metrics{ID: "2", MType: domain.Counter, Value: nil}

	mockRanger.On("Range", mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(0).(func(string, *domain.Metrics) bool)
		callback("gauge_1", metric1)
		callback("counter_2", metric2)
	}).Return()

	result := getAllService.GetAll()

	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Contains(t, result, metric1)
	assert.Contains(t, result, metric2)
	mockRanger.AssertExpectations(t)
}

func TestUpdateMetricService_Update_Default(t *testing.T) {
	// Моки для UpdateGaugeMetricService и UpdateCounterMetricService
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)
	mockEncoder := new(MockKeyEncoder)

	gaugeService := NewUpdateGaugeMetricService(mockSaver, mockGetter, mockEncoder)
	counterService := NewUpdateCounterMetricService(mockSaver, mockGetter, mockEncoder)

	// Новый фасадный сервис для метрик
	updateService := NewUpdateMetricService(gaugeService, counterService)

	// Метрика с неподдерживаемым типом
	invalidMetric := &domain.Metrics{
		ID:    "1",
		MType: "InvalidType", // Указан неподдерживаемый тип
	}

	// Вызываем метод Update для неподдерживаемого типа
	result := updateService.UpdateMetric(invalidMetric)

	// Проверяем, что результат равен nil
	assert.Nil(t, result)

	// Проверяем, что методы сохранения и получения не были вызваны
	mockSaver.AssertNotCalled(t, "Save")
	mockGetter.AssertNotCalled(t, "Get")
	mockEncoder.AssertNotCalled(t, "Encode")
}
