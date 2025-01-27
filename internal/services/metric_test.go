package services

import (
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/types"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorageEngine mocks the StorageEngineInterface
type MockStorageEngine struct {
	mock.Mock
}

func (m *MockStorageEngine) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockStorageEngine) Set(key, value string) {
	m.Called(key, value)
}

func (m *MockStorageEngine) Generate() <-chan [2]string {
	ch := make(chan [2]string, 1)
	ch <- [2]string{"gauge:cpu", "42.5"}
	close(ch)
	return ch
}

// MockKeyEngine mocks the KeyEngineInterface
type MockKeyEngine struct {
	mock.Mock
}

func (m *MockKeyEngine) Encode(metricType, name string) string {
	args := m.Called(metricType, name)
	return args.String(0)
}

func (m *MockKeyEngine) Decode(key string) (string, string, error) {
	args := m.Called(key)
	return args.String(0), args.String(1), args.Error(2)
}

// MockStrategyEngine mocks the StrategyEngineInterface
type MockStrategyEngine struct {
	mock.Mock
}

func (m *MockStrategyEngine) Update(currentValue, newValue string) (string, error) {
	args := m.Called(currentValue, newValue)
	return args.String(0), args.Error(1)
}

// Test UpdateMetric - Valid Case
func TestMetricService_UpdateMetric(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	mockStrategy := new(MockStrategyEngine)

	strategyEngines := map[types.MetricType]engines.StrategyUpdateEngineInterface{
		types.GaugeType: mockStrategy,
	}

	service := NewMetricService(mockStorage, strategyEngines, mockKeyEngine)

	req := types.UpdateMetricRequest{
		Type:  "gauge",
		Name:  "cpu",
		Value: "50.5",
	}

	mockKeyEngine.On("Encode", "gauge", "cpu").Return("gauge:cpu")
	mockStorage.On("Get", "gauge:cpu").Return("42.5", true)
	mockStrategy.On("Update", "42.5", "50.5").Return("50.5", nil)
	mockStorage.On("Set", "gauge:cpu", "50.5").Return()

	err := service.UpdateMetric(req)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
	mockStrategy.AssertExpectations(t)
}

// Test UpdateMetric - Invalid Metric Type
func TestMetricService_UpdateMetric_InvalidType(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	mockStrategy := new(MockStrategyEngine)

	strategyEngines := map[types.MetricType]engines.StrategyUpdateEngineInterface{
		types.GaugeType: mockStrategy,
	}

	service := NewMetricService(mockStorage, strategyEngines, mockKeyEngine)

	req := types.UpdateMetricRequest{
		Type:  "", // Invalid type
		Name:  "cpu",
		Value: "50.5",
	}

	err := service.UpdateMetric(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
	assert.Equal(t, "metric type is required", err.Error())
}

// Test GetMetric - Valid Case
func TestMetricService_GetMetric(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)

	service := NewMetricService(mockStorage, nil, mockKeyEngine)

	req := types.GetMetricRequest{
		Type: "gauge",
		Name: "cpu",
	}

	mockKeyEngine.On("Encode", "gauge", "cpu").Return("gauge:cpu")
	mockStorage.On("Get", "gauge:cpu").Return("50.5", true)

	value, err := service.GetMetric(req)

	assert.NoError(t, err)
	assert.Equal(t, "50.5", value)
	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

// Test GetMetric - Metric Not Found
func TestMetricService_GetMetric_NotFound(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)

	service := NewMetricService(mockStorage, nil, mockKeyEngine)

	req := types.GetMetricRequest{
		Type: "gauge",
		Name: "memory",
	}

	mockKeyEngine.On("Encode", "gauge", "memory").Return("gauge:memory")
	mockStorage.On("Get", "gauge:memory").Return("", false)

	value, err := service.GetMetric(req)

	assert.Error(t, err)
	assert.Equal(t, "", value)
	assert.Equal(t, http.StatusNotFound, err.Status())
	assert.Equal(t, "metric not found", err.Error())

	mockStorage.AssertExpectations(t)
	mockKeyEngine.AssertExpectations(t)
}

// Test GetAllMetrics
func TestMetricService_GetAllMetrics(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)

	service := NewMetricService(mockStorage, nil, mockKeyEngine)

	mockKeyEngine.On("Decode", "gauge:cpu").Return("gauge", "cpu", nil)

	metrics := service.GetAllMetrics()

	assert.Len(t, metrics, 1)
	assert.Equal(t, "cpu", metrics[0].Name)
	assert.Equal(t, "42.5", metrics[0].Value)
	assert.Equal(t, "gauge", metrics[0].Type)

	mockKeyEngine.AssertExpectations(t)
}

// Test UpdateMetric - Invalid Metric Name
func TestMetricService_UpdateMetric_InvalidName(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	mockStrategy := new(MockStrategyEngine)

	strategyEngines := map[types.MetricType]engines.StrategyUpdateEngineInterface{
		types.GaugeType: mockStrategy,
	}

	service := NewMetricService(mockStorage, strategyEngines, mockKeyEngine)

	req := types.UpdateMetricRequest{
		Type:  "gauge",
		Name:  "", // Invalid name
		Value: "50.5",
	}

	err := service.UpdateMetric(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
	assert.Equal(t, "metric name is required", err.Error())
}

// Test UpdateMetric - Invalid Metric Value
func TestMetricService_UpdateMetric_InvalidValue(t *testing.T) {
	mockStorage := new(MockStorageEngine)
	mockKeyEngine := new(MockKeyEngine)
	mockStrategy := new(MockStrategyEngine)

	strategyEngines := map[types.MetricType]engines.StrategyUpdateEngineInterface{
		types.GaugeType: mockStrategy,
	}

	service := NewMetricService(mockStorage, strategyEngines, mockKeyEngine)

	req := types.UpdateMetricRequest{
		Type:  "gauge",
		Name:  "cpu",
		Value: "", // Invalid value
	}

	err := service.UpdateMetric(req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
	assert.Equal(t, "metric value is required", err.Error())
}
