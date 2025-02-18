package services

import (
	"go-metrics-alerting/internal/types"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGaugeSaver struct {
	mock.Mock
}

func (m *MockGaugeSaver) Save(key string, value float64) bool {
	args := m.Called(key, value)
	return args.Bool(0)
}

type MockGaugeGetter struct {
	mock.Mock
}

func (m *MockGaugeGetter) Get(key string) (float64, bool) {
	args := m.Called(key)
	return args.Get(0).(float64), args.Bool(1)
}

type MockCounterSaver struct {
	mock.Mock
}

func (m *MockCounterSaver) Save(key string, value int64) bool {
	args := m.Called(key, value)
	return args.Bool(0)
}

type MockCounterGetter struct {
	mock.Mock
}

func (m *MockCounterGetter) Get(key string) (int64, bool) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Bool(1)
}

func TestUpdateMetricsService_Update_Gauge(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: types.Gauge,
		Value: float64Ptr(5.5),
	}

	// Expect Save to be called once for the gauge
	mockGaugeSaver.On("Save", "metric1", 5.5).Return(true)

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	assert.Equal(t, float64(5.5), *resp.Value)
	mockGaugeSaver.AssertExpectations(t)
}

func TestUpdateMetricsService_Update_Counter_ExistingValue(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: types.Counter,
		Delta: int64Ptr(3),
	}

	// Setup mock: counter exists with value 2
	mockCounterGetter.On("Get", "metric1").Return(int64(2), true)
	// Expect Save to be called with the updated value: 2 + 3 = 5
	mockCounterSaver.On("Save", "metric1", int64(5)).Return(true)

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	mockCounterGetter.AssertExpectations(t)
	mockCounterSaver.AssertExpectations(t)
}

func TestUpdateMetricsService_Update_Counter_NoExistingValue(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: types.Counter,
		Delta: int64Ptr(3),
	}

	// Setup mock: counter does not exist, so return default value 0
	mockCounterGetter.On("Get", "metric1").Return(int64(0), false)
	// Expect Save to be called with the updated value: 0 + 3 = 3
	mockCounterSaver.On("Save", "metric1", int64(3)).Return(true)

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	mockCounterGetter.AssertExpectations(t)
	mockCounterSaver.AssertExpectations(t)
}

func TestUpdateMetricsService_Update_UnsupportedType(t *testing.T) {
	mockGaugeSaver := new(MockGaugeSaver)
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterSaver := new(MockCounterSaver)
	mockCounterGetter := new(MockCounterGetter)

	// Create the service
	service := NewUpdateMetricsService(mockGaugeSaver, mockGaugeGetter, mockCounterSaver, mockCounterGetter)

	req := &types.UpdateMetricsRequest{
		ID:    "metric1",
		MType: "Unsupported", // This is an unsupported type
	}

	// Call the Update method
	resp, err := service.Update(req)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, resp) // No response for unsupported types
}

// Тест для успешного получения значения Gauge
func TestGetMetricValueService_GetMetricValue_Gauge(t *testing.T) {
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterGetter := new(MockCounterGetter)

	// Создание сервиса
	service := NewGetMetricValueService(mockGaugeGetter, mockCounterGetter)

	// Настроим ожидания для мока
	mockGaugeGetter.On("Get", "metric1").Return(5.5, true)

	req := &types.GetMetricValueRequest{
		ID:    "metric1",
		MType: types.Gauge,
	}

	// Вызов метода
	resp, err := service.GetMetricValue(req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	assert.Equal(t, "5.5", resp.Value)
	mockGaugeGetter.AssertExpectations(t)
}

// Тест для успешного получения значения Counter
func TestGetMetricValueService_GetMetricValue_Counter(t *testing.T) {
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterGetter := new(MockCounterGetter)

	// Создание сервиса
	service := NewGetMetricValueService(mockGaugeGetter, mockCounterGetter)

	// Настроим ожидания для мока
	mockCounterGetter.On("Get", "metric1").Return(int64(10), true)

	req := &types.GetMetricValueRequest{
		ID:    "metric1",
		MType: types.Counter,
	}

	// Вызов метода
	resp, err := service.GetMetricValue(req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metric1", resp.ID)
	assert.Equal(t, "10", resp.Value)
	mockCounterGetter.AssertExpectations(t)
}

// Тест для случая, когда метрика не найдена (Gauge)
func TestGetMetricValueService_GetMetricValue_Gauge_NotFound(t *testing.T) {
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterGetter := new(MockCounterGetter)

	// Создание сервиса
	service := NewGetMetricValueService(mockGaugeGetter, mockCounterGetter)

	// Настроим ожидания для мока
	mockGaugeGetter.On("Get", "metric1").Return(0.0, false)

	req := &types.GetMetricValueRequest{
		ID:    "metric1",
		MType: types.Gauge,
	}

	// Вызов метода
	resp, err := service.GetMetricValue(req)

	// Проверка
	assert.NoError(t, err)
	assert.Nil(t, resp)
	mockGaugeGetter.AssertExpectations(t)
}

// Тест для случая, когда метрика не найдена (Counter)
func TestGetMetricValueService_GetMetricValue_Counter_NotFound(t *testing.T) {
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterGetter := new(MockCounterGetter)

	// Создание сервиса
	service := NewGetMetricValueService(mockGaugeGetter, mockCounterGetter)

	// Настроим ожидания для мока
	mockCounterGetter.On("Get", "metric1").Return(int64(0), false)

	req := &types.GetMetricValueRequest{
		ID:    "metric1",
		MType: types.Counter,
	}

	// Вызов метода
	resp, err := service.GetMetricValue(req)

	// Проверка
	assert.NoError(t, err)
	assert.Nil(t, resp)
	mockCounterGetter.AssertExpectations(t)
}

// Тест для случая, когда тип метрики не поддерживается
func TestGetMetricValueService_GetMetricValue_UnsupportedType(t *testing.T) {
	mockGaugeGetter := new(MockGaugeGetter)
	mockCounterGetter := new(MockCounterGetter)

	// Создание сервиса
	service := NewGetMetricValueService(mockGaugeGetter, mockCounterGetter)

	req := &types.GetMetricValueRequest{
		ID:    "metric1",
		MType: "Unsupported", // Не поддерживаемый тип
	}

	// Вызов метода
	resp, err := service.GetMetricValue(req)

	// Проверка
	assert.NoError(t, err)
	assert.Nil(t, resp)
}

// Мок для Ranger, реализующий метод Range для Gauge
type MockGaugeRanger struct {
	mock.Mock
}

func (m *MockGaugeRanger) Range(f func(key string, value float64) bool) {
	args := m.Called(f)
	if args.Bool(0) {
		f("metric1", 5.5)
		f("metric2", 10.5)
	}
}

// Мок для Ranger, реализующий метод Range для Counter
type MockCounterRanger struct {
	mock.Mock
}

func (m *MockCounterRanger) Range(f func(key string, value int64) bool) {
	args := m.Called(f)
	if args.Bool(0) {
		f("metric3", int64(100))
	}
}

// Тест для получения всех значений метрик
func TestGetAllMetricValuesService_GetAllMetricValues(t *testing.T) {
	mockGaugeRanger := new(MockGaugeRanger)
	mockCounterRanger := new(MockCounterRanger)

	// Создание сервиса
	service := NewGetAllMetricValuesService(mockGaugeRanger, mockCounterRanger)

	// Настроим ожидания для мока
	mockGaugeRanger.On("Range", mock.Anything).Return(true)
	mockCounterRanger.On("Range", mock.Anything).Return(true)

	// Вызов метода
	resp := service.GetAllMetricValues()

	// Проверка
	assert.Len(t, resp, 3)
	assert.Equal(t, "metric1", resp[0].ID)
	assert.Equal(t, "5.5", resp[0].Value)
	assert.Equal(t, "metric2", resp[1].ID)
	assert.Equal(t, "10.5", resp[1].Value)
	assert.Equal(t, "metric3", resp[2].ID)
	assert.Equal(t, "100", resp[2].Value)

	mockGaugeRanger.AssertExpectations(t)
	mockCounterRanger.AssertExpectations(t)
}

func float64Ptr(v float64) *float64 {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}
