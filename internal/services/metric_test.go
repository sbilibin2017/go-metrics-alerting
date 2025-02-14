package services

import (
	"errors"
	"go-metrics-alerting/internal/types"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricRepository имитирует поведение хранилища метрик.
type MockMetricRepository struct {
	mock.Mock
}

func (m *MockMetricRepository) Save(metricType, metricName, metricValue string) {
	m.Called(metricType, metricName, metricValue)
}

func (m *MockMetricRepository) Get(metricType, metricName string) (string, error) {
	args := m.Called(metricType, metricName)
	return args.String(0), args.Error(1)
}

func (m *MockMetricRepository) GetAll() [][]string {
	args := m.Called()
	return args.Get(0).([][]string)
}

func TestUpdateMetric_Success_GaugeType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)

	req := &types.UpdateMetricValueRequest{
		Name:  "metric1",
		Type:  types.Gauge,
		Value: "10.5", // New value to replace the current value
	}

	// Mock the Get method to return an existing value for the metric (e.g., 5.5)
	mockRepo.On("Get", string(types.Gauge), "metric1").Return("5.5", nil).Once()

	// Mock the Save method to save the new value (replacing the old one)
	mockRepo.On("Save", string(types.Gauge), "metric1", "10.5").Return().Once()

	// Call the UpdateMetric method
	service.UpdateMetric(req)

	// Assert that the repository expectations are met (i.e., Get and Save were called correctly)
	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_RepoReturnsEmpty_FallsBackToDefault(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)

	req := &types.UpdateMetricValueRequest{
		Name:  "metric1",
		Type:  types.Counter,
		Value: "10",
	}

	// Mock the Get method to return an empty string (MetricEmptyString)
	mockRepo.On("Get", string(types.Counter), "metric1").Return(MetricEmptyString, nil).Once()

	// Mock the Save method to do nothing (i.e., no error returned)
	mockRepo.On("Save", string(types.Counter), "metric1", "10").Return().Once()

	// Call the UpdateMetric method
	service.UpdateMetric(req)

	// Assert that the fallback logic worked (the value should be set to DefaultMetricValue)
	mockRepo.AssertExpectations(t)
}

func TestNewMetricService(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestGetMetric_Success(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)

	mockRepo.On("Get", string(types.Counter), "metric1").Return("10", nil)

	value, err := service.GetMetric(&types.GetMetricValueRequest{
		Name: "metric1",
		Type: types.Counter,
	})

	// Check for no error
	assert.Nil(t, err)
	// Assert that the value returned is correct
	assert.Equal(t, "10", value)
	mockRepo.AssertExpectations(t)
}

func TestGetMetric_Error(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)

	// Simulate an error scenario
	mockRepo.On("Get", string(types.Counter), "metric1").Return("", errors.New("not found"))

	// Call the GetMetric method
	value, errResp := service.GetMetric(&types.GetMetricValueRequest{
		Name: "metric1",
		Type: types.Counter,
	})

	// Check that an error response is returned
	assert.Equal(t, "", value)
	assert.NotNil(t, errResp)
	assert.Equal(t, http.StatusNotFound, errResp.Code)
	assert.Equal(t, "value not found", errResp.Message)

	mockRepo.AssertExpectations(t)
}

func TestListMetrics_Success(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)

	// Mock GetAll to return a list of metrics
	mockRepo.On("GetAll").Return([][]string{
		{"counter", "metric1", "10"},
		{"gauge", "metric2", "3.14"},
	})

	metrics := service.ListMetrics()

	// Assert that the correct number of metrics is returned
	assert.Len(t, metrics, 2)
	// Assert that the metric names and values are correct
	assert.Equal(t, "metric1", metrics[0].Name)
	assert.Equal(t, "10", metrics[0].Value)
	assert.Equal(t, "metric2", metrics[1].Name)
	assert.Equal(t, "3.14", metrics[1].Value)

	mockRepo.AssertExpectations(t)
}

func TestListMetrics_Empty(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := NewMetricService(mockRepo)

	// Mock GetAll to return an empty list of metrics
	mockRepo.On("GetAll").Return([][]string{})

	metrics := service.ListMetrics()

	// Assert that the returned metrics list is empty
	assert.Len(t, metrics, 0)
	mockRepo.AssertExpectations(t)
}
