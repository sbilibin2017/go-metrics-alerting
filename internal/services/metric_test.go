package services

import (
	"fmt"
	"go-metrics-alerting/internal/types"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the dependencies

type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

type MockMetricUpdateStrategy struct {
	mock.Mock
}

func (m *MockMetricUpdateStrategy) Update(req *types.MetricsRequest, currentValue string) (*types.MetricsRequest, error) {
	args := m.Called(req, currentValue)
	return args.Get(0).(*types.MetricsRequest), args.Error(1)
}

type MockIDValidator struct {
	mock.Mock
}

func (m *MockIDValidator) Validate(id string) bool {
	args := m.Called(id)
	return args.Bool(0)
}

type MockMTypeValidator struct {
	mock.Mock
}

func (m *MockMTypeValidator) Validate(mType string) bool {
	args := m.Called(mType)
	return args.Bool(0)
}

type MockDeltaValidator struct {
	mock.Mock
}

func (m *MockDeltaValidator) Validate(mtype string, delta *int64) bool {
	args := m.Called(mtype, delta)
	return args.Bool(0)
}

type MockValueValidator struct {
	mock.Mock
}

func (m *MockValueValidator) Validate(mtype string, value *float64) bool {
	args := m.Called(mtype, value)
	return args.Bool(0)
}

func TestUpdateMetric_Fail_InvalidID(t *testing.T) {
	// Setup
	getter := new(MockGetter)
	saver := new(MockSaver)
	strategy := new(MockMetricUpdateStrategy)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	idValidator.On("Validate", "metric2").Return(false)

	service := NewUpdateMetricService(getter, saver, map[string]MetricUpdateStrategy{"counter": strategy}, idValidator, mtypeValidator, deltaValidator, valueValidator)

	// Execution
	req := &types.MetricsRequest{
		ID:    "metric2",
		MType: "counter",
		Delta: ptrInt64(10),
	}
	resp, err := service.UpdateMetric(req)

	// Validation
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusNotFound, err.Status)
	assert.Equal(t, "Metric is not found", err.Message)
}

func TestUpdateMetric_Fail_InvalidMetricType(t *testing.T) {
	// Setup
	getter := new(MockGetter)
	saver := new(MockSaver)
	strategy := new(MockMetricUpdateStrategy)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	idValidator.On("Validate", "metric3").Return(true)
	mtypeValidator.On("Validate", "invalidType").Return(false)

	service := NewUpdateMetricService(getter, saver, map[string]MetricUpdateStrategy{"counter": strategy}, idValidator, mtypeValidator, deltaValidator, valueValidator)

	// Execution
	req := &types.MetricsRequest{
		ID:    "metric3",
		MType: "invalidType",
		Delta: ptrInt64(10),
	}
	resp, err := service.UpdateMetric(req)

	// Validation
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, "Invalid metric type", err.Message)
}

func TestUpdateMetric_Fail_InvalidValueForGauge(t *testing.T) {
	// Setup
	getter := new(MockGetter)
	saver := new(MockSaver)
	strategy := new(MockMetricUpdateStrategy)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	// Mock Setup
	idValidator.On("Validate", "metric5").Return(true)
	mtypeValidator.On("Validate", "gauge").Return(true)
	deltaValidator.On("Validate", "gauge", mock.AnythingOfType("*int64")).Return(true)
	// Here we fix the mock expectation for the nil value to be *float64(nil)
	valueValidator.On("Validate", "gauge", (*float64)(nil)).Return(false)

	service := NewUpdateMetricService(getter, saver, map[string]MetricUpdateStrategy{"gauge": strategy}, idValidator, mtypeValidator, deltaValidator, valueValidator)

	// Execution
	req := &types.MetricsRequest{
		ID:    "metric5",
		MType: "gauge",
		Value: nil, // Sending nil for the gauge value
	}
	resp, err := service.UpdateMetric(req)

	// Validation
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, "Invalid value for Gauge metric", err.Message)
}

func TestUpdateMetric_Fail_InvalidDeltaForCounter(t *testing.T) {
	// Setup
	getter := new(MockGetter)
	saver := new(MockSaver)
	strategy := new(MockMetricUpdateStrategy)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	// Mock Setup
	idValidator.On("Validate", "metric6").Return(true)
	mtypeValidator.On("Validate", "counter").Return(true)
	// Mock the DeltaValidator to return false when Delta is invalid for Counter metrics
	deltaValidator.On("Validate", "counter", mock.AnythingOfType("*int64")).Return(false)
	// Mock valueValidator to return true (not needed for this test, but required for execution)
	valueValidator.On("Validate", "counter", mock.AnythingOfType("*float64")).Return(true)

	service := NewUpdateMetricService(getter, saver, map[string]MetricUpdateStrategy{"counter": strategy}, idValidator, mtypeValidator, deltaValidator, valueValidator)

	// Execution
	delta := int64(100) // example of a delta value
	req := &types.MetricsRequest{
		ID:    "metric6",
		MType: "counter",
		Delta: &delta, // setting delta for the counter metric
	}
	resp, err := service.UpdateMetric(req)

	// Validation
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, "Invalid delta for Counter metric", err.Message)
}

func TestUpdateMetric_MetricNotFoundInStorage(t *testing.T) {
	// Setup
	getter := new(MockGetter)
	saver := new(MockSaver)
	strategy := new(MockMetricUpdateStrategy)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	// Mock Setup
	idValidator.On("Validate", "metric7").Return(true)
	mtypeValidator.On("Validate", "gauge").Return(true)
	deltaValidator.On("Validate", "gauge", mock.AnythingOfType("*int64")).Return(true)
	valueValidator.On("Validate", "gauge", mock.AnythingOfType("*float64")).Return(true)

	// Mock the getter to return an error when fetching the current value (simulating metric not found)
	getter.On("Get", "metric7").Return("", fmt.Errorf("metric not found"))

	// Mock the saver to just return nil (we'll focus on the getter in this test)
	saver.On("Save", "metric7", "gauge").Return(nil)

	// Execution
	value := float64(12.5)

	// Mock the strategy's Update method to return the updated request (just return the input request in this case)
	strategy.On("Update", mock.AnythingOfType("*types.MetricsRequest"), "0").Return(&types.MetricsRequest{
		ID:    "metric7",
		MType: "gauge",
		Delta: nil,
		Value: &value, // ensure the value passed here is what we expect
	}, nil)

	service := NewUpdateMetricService(getter, saver, map[string]MetricUpdateStrategy{"gauge": strategy}, idValidator, mtypeValidator, deltaValidator, valueValidator)

	// example of a value for gauge
	req := &types.MetricsRequest{
		ID:    "metric7",
		MType: "gauge",
		Value: &value, // setting value for the gauge metric
	}
	resp, _ := service.UpdateMetric(req)

	// Validation

	assert.Equal(t, "metric7", resp.MetricsRequest.ID)
	assert.Equal(t, "gauge", resp.MetricsRequest.MType)
	assert.Equal(t, &value, resp.MetricsRequest.Value) // value should be updated as expected

	// Ensure the current value was assumed to be "0"
	getter.AssertExpectations(t)
	strategy.AssertExpectations(t) // Add this to ensure the strategy's Update method was called as expected
}

func TestUpdateMetric_NoStrategyFound(t *testing.T) {
	// Setup mocks
	getter := new(MockGetter)
	saver := new(MockSaver)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)

	// Mock Setup
	// Let's assume the ID is valid
	idValidator.On("Validate", "metric1").Return(true)
	// The metric type "unknown" is provided, and we expect that it doesn't have a strategy
	mtypeValidator.On("Validate", "unknown").Return(true)
	// Delta and Value validators can return true as their logic isn't the focus here
	deltaValidator.On("Validate", "unknown", mock.AnythingOfType("*int64")).Return(true)
	valueValidator.On("Validate", "unknown", mock.AnythingOfType("*float64")).Return(true)

	// Mock the getter to return a valid metric value
	getter.On("Get", "metric1").Return("0", nil)

	// Mock the saver to just return nil (no saving needed in this case)
	saver.On("Save", "metric1", "unknown").Return(nil)

	// Create the service with an empty strategy map
	service := NewUpdateMetricService(
		getter,
		saver,
		map[string]MetricUpdateStrategy{}, // Empty strategy map to simulate no strategy for "unknown"
		idValidator,
		mtypeValidator,
		deltaValidator,
		valueValidator,
	)

	// Create a request with an invalid MType that doesn't have a strategy
	req := &types.MetricsRequest{
		ID:    "metric1",
		MType: "unknown", // This type doesn't have any strategy
	}

	// Execution
	resp, err := service.UpdateMetric(req)

	// Validation
	// Expect an error and the message "No strategy found for metric type"
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, "No strategy found for metric type", err.Message)

}

type MockStringGetter struct {
	mock.Mock
}

func (m *MockStringGetter) Get(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

type MockStringSaver struct {
	mock.Mock
}

func (m *MockStringSaver) Save(id, mtype string) error {
	args := m.Called(id, mtype)
	return args.Error(0)
}

type MockMetricUpdateStrategy2 struct {
	mock.Mock
}

func (m *MockMetricUpdateStrategy2) Update(req *types.MetricsRequest, currentValue string) (*types.MetricsRequest, error) {
	args := m.Called(req, currentValue)
	return nil, args.Error(1) // Simulate error in strategy update
}

func TestUpdateMetric_StrategyUpdateError(t *testing.T) {
	// Mocks Setup
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)
	stringGetter := new(MockStringGetter)
	stringSaver := new(MockStringSaver)

	// Define mock behavior for each component
	idValidator.On("Validate", "metric8").Return(true)
	mtypeValidator.On("Validate", "counter").Return(true)
	deltaValidator.On("Validate", "counter", mock.Anything).Return(true)

	// Adjust mock for valueValidator: since this is a "counter", value should be nil (we pass a nil pointer here)
	valueValidator.On("Validate", "counter", mock.Anything).Return(true) // For counter, we expect no value validation

	// Mock the string getter to return a value
	stringGetter.On("Get", "metric8").Return("5", nil)

	// Create the mock strategy that simulates an error during Update
	strategy := new(MockMetricUpdateStrategy2)
	strategy.On("Update", mock.Anything, "5").Return(nil, fmt.Errorf("strategy update error"))

	// Create the service
	service := NewUpdateMetricService(
		stringGetter,
		stringSaver,
		map[string]MetricUpdateStrategy{"counter": strategy},
		idValidator,
		mtypeValidator,
		deltaValidator,
		valueValidator,
	)

	// Define the request
	delta := int64(10)
	req := &types.MetricsRequest{
		ID:    "metric8",
		MType: "counter",
		Delta: &delta,
	}

	// Call UpdateMetric method
	resp, errResp := service.UpdateMetric(req)

	// Assertions
	// Ensure response is nil and error response is correctly set
	assert.Nil(t, resp)
	assert.NotNil(t, errResp)
	assert.Equal(t, http.StatusBadRequest, errResp.Status)
	assert.Equal(t, "Metric is not updated", errResp.Message)

	// Check that the correct expectations were met for the mock strategy
	strategy.AssertExpectations(t)
}

func TestUpdateMetric_SaveError(t *testing.T) {
	// Mocks Setup
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	deltaValidator := new(MockDeltaValidator)
	valueValidator := new(MockValueValidator)
	stringGetter := new(MockStringGetter)
	stringSaver := new(MockStringSaver)

	// Define mock behavior for each component
	idValidator.On("Validate", "metric8").Return(true)
	mtypeValidator.On("Validate", "counter").Return(true)
	deltaValidator.On("Validate", "counter", mock.Anything).Return(true)

	// Adjust mock for valueValidator: since this is a "counter", value should be nil (we pass a nil pointer here)
	valueValidator.On("Validate", "counter", mock.Anything).Return(true)

	// Mock the string getter to return a value
	stringGetter.On("Get", "metric8").Return("5", nil)

	// Create the mock strategy that simulates success during Update
	strategy := new(MockMetricUpdateStrategy)
	strategy.On("Update", mock.Anything, "5").Return(&types.MetricsRequest{
		ID:    "metric8",
		MType: "counter",
		Delta: nil,
		Value: nil,
	}, nil) // Simulating successful update of the metric

	// Simulate error during saving the metric
	stringSaver.On("Save", "metric8", "counter").Return(fmt.Errorf("save error"))

	// Create the service
	service := NewUpdateMetricService(
		stringGetter,
		stringSaver,
		map[string]MetricUpdateStrategy{"counter": strategy},
		idValidator,
		mtypeValidator,
		deltaValidator,
		valueValidator,
	)

	// Define the request
	delta := int64(10)
	req := &types.MetricsRequest{
		ID:    "metric8",
		MType: "counter",
		Delta: &delta,
	}

	// Call UpdateMetric method
	resp, errResp := service.UpdateMetric(req)

	// Assertions
	// Ensure response is nil and error response is correctly set for saving failure
	assert.Nil(t, resp)
	assert.NotNil(t, errResp)
	assert.Equal(t, http.StatusInternalServerError, errResp.Status)
	assert.Equal(t, "Metric is not saved", errResp.Message)

	// Check that the correct expectations were met for the mock strategy and save operation
	stringSaver.AssertExpectations(t)
	strategy.AssertExpectations(t)
}

func TestGetMetric_Success(t *testing.T) {
	// Setup
	getter := new(MockStringGetter)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)

	// Мокаем поведение
	idValidator.On("Validate", "metric1").Return(true)
	mtypeValidator.On("Validate", "gauge").Return(true)
	getter.On("Get", "metric1").Return("123.45", nil)

	// Создаём сервис
	service := NewGetMetricService(getter, idValidator, mtypeValidator)

	// Запрос
	req := &types.MetricValueRequest{
		ID:    "metric1",
		MType: "gauge",
	}
	value, resp := service.GetMetric(req)

	// Проверки
	assert.NotNil(t, value)
	assert.Equal(t, "123.45", *value)
	assert.Nil(t, resp) // Не должно быть ошибки
}

func TestGetMetric_Fail_InvalidID(t *testing.T) {
	// Setup
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)
	stringGetter := new(MockStringGetter)

	// Setup mock behavior
	idValidator.On("Validate", "invalidMetricID").Return(false)

	// Create service instance
	service := NewGetMetricService(stringGetter, idValidator, mtypeValidator)

	// Execution
	req := &types.MetricValueRequest{
		ID:    "invalidMetricID",
		MType: "gauge", // assume the type is valid for this test
	}
	metricValue, resp := service.GetMetric(req)

	// Validation
	assert.Nil(t, metricValue) // Expecting nil response
	assert.Equal(t, http.StatusNotFound, resp.Status)
	assert.Equal(t, "Metric is not found", resp.Message)
}

func TestGetMetric_Fail_InvalidMetricType(t *testing.T) {
	// Setup
	getter := new(MockStringGetter)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)

	// Mock behavior
	idValidator.On("Validate", "metric1").Return(true)
	mtypeValidator.On("Validate", "invalidType").Return(false)

	// Create service instance
	service := NewGetMetricService(getter, idValidator, mtypeValidator)

	// Execution
	req := &types.MetricValueRequest{
		ID:    "metric1",
		MType: "invalidType",
	}
	_, resp := service.GetMetric(req)

	// Validation
	assert.NotNil(t, resp) // Ensure the response is not nil
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, "Invalid metric type", resp.Message)
}

func TestGetMetric_MetricNotFound(t *testing.T) {
	// Setup
	getter := new(MockStringGetter)
	idValidator := new(MockIDValidator)
	mtypeValidator := new(MockMTypeValidator)

	// Mock behavior
	idValidator.On("Validate", "metric1").Return(true)
	mtypeValidator.On("Validate", "gauge").Return(true)
	getter.On("Get", "metric1").Return("", fmt.Errorf("metric not found"))

	// Create service instance
	service := NewGetMetricService(getter, idValidator, mtypeValidator)

	// Execution
	req := &types.MetricValueRequest{
		ID:    "metric1",
		MType: "gauge",
	}
	_, resp := service.GetMetric(req)

	// Validation
	assert.NotNil(t, resp) // Ensure the response is not nil
	assert.Equal(t, http.StatusNotFound, resp.Status)
	assert.Equal(t, "Metric is not found", resp.Message)
}

// Define a mock for the Ranger interface
type MockRanger struct {
	mock.Mock
}

func (m *MockRanger) Range(f func(key string, value string) bool) {
	args := m.Called(f)
	// Call the function for each key-value pair, simulating Ranger behavior
	for _, arg := range args.Get(0).([][2]string) {
		f(arg[0], arg[1])
	}
}

func TestListMetrics_EmptyRanger(t *testing.T) {
	// Setup
	mockRanger := new(MockRanger)
	mockRanger.On("Range", mock.Anything).Return([][2]string{})

	service := NewListMetricsService(mockRanger)

	// Execution
	metrics := service.ListMetrics()

	// Validation
	assert.Empty(t, metrics)
}

func TestListMetrics_SingleMetric(t *testing.T) {
	// Setup
	mockRanger := new(MockRanger)
	mockRanger.On("Range", mock.Anything).Return([][2]string{
		{"metric1", "10"},
	})

	service := NewListMetricsService(mockRanger)

	// Execution
	metrics := service.ListMetrics()

	// Validation
	assert.Len(t, metrics, 1)
	assert.Equal(t, "metric1", metrics[0].ID)
	assert.Equal(t, "10", metrics[0].Value)
}

func TestListMetrics_MultipleMetrics(t *testing.T) {
	// Setup
	mockRanger := new(MockRanger)
	mockRanger.On("Range", mock.Anything).Return([][2]string{
		{"metric1", "10"},
		{"metric2", "20"},
		{"metric3", "30"},
	})

	service := NewListMetricsService(mockRanger)

	// Execution
	metrics := service.ListMetrics()

	// Validation
	assert.Len(t, metrics, 3)
	assert.Equal(t, "metric1", metrics[0].ID)
	assert.Equal(t, "10", metrics[0].Value)
	assert.Equal(t, "metric2", metrics[1].ID)
	assert.Equal(t, "20", metrics[1].Value)
	assert.Equal(t, "metric3", metrics[2].ID)
	assert.Equal(t, "30", metrics[2].Value)
}

func ptrInt64(v int64) *int64 {
	return &v
}
