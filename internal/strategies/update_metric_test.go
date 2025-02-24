package strategies_test

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/strategies"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(key string, value *domain.Metrics) {
	m.Called(key, value)
}

type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) (*domain.Metrics, bool) {
	args := m.Called(key)
	if args.Get(0) == nil {
		// Return nil for *domain.Metrics and false for existence flag
		return nil, false
	}
	return args.Get(0).(*domain.Metrics), args.Bool(1)
}

func TestUpdateGaugeMetricStrategy(t *testing.T) {
	// Mock instances
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)

	// Create the strategy
	strategy := strategies.NewUpdateGaugeMetricStrategy(mockSaver, mockGetter)

	// Define a metric to update
	metric := &domain.Metrics{
		ID:    "test_metric",
		MType: domain.Gauge,
	}

	// Test saving the metric
	mockSaver.On("Save", "test_metric:gauge", metric).Once()

	// Call the method
	updatedMetric := strategy.UpdateMetric(metric)

	// Assert that the metric returned is the same as the one passed
	assert.Equal(t, metric, updatedMetric)

	// Assert the expected method call
	mockSaver.AssertExpectations(t)
}

func TestUpdateCounterMetricStrategy_NewMetric(t *testing.T) {
	// Mock instances
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)

	// Create the strategy
	strategy := strategies.NewUpdateCounterMetricStrategy(mockSaver, mockGetter)

	// Define a new metric to update
	metric := &domain.Metrics{
		ID:    "test_metric",
		MType: domain.Counter,
		Delta: new(int64),
	}

	// Mock the Get method to return nothing for the new metric
	mockGetter.On("Get", "test_metric:counter").Return(nil, false)

	// Test saving the metric when it's new
	mockSaver.On("Save", "test_metric:counter", metric).Once()

	// Call the method
	updatedMetric := strategy.UpdateMetric(metric)

	// Assert that the metric returned is the same as the one passed
	assert.Equal(t, metric, updatedMetric)

	// Assert the expected method calls
	mockSaver.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
}

func TestUpdateCounterMetricStrategy_ExistingMetric(t *testing.T) {
	// Mock instances
	mockSaver := new(MockSaver)
	mockGetter := new(MockGetter)

	// Create the strategy
	strategy := strategies.NewUpdateCounterMetricStrategy(mockSaver, mockGetter)

	// Define a new metric and an existing metric
	metric := &domain.Metrics{
		ID:    "test_metric",
		MType: domain.Counter,
		Delta: new(int64),
	}
	existingMetric := &domain.Metrics{
		ID:    "test_metric",
		MType: domain.Counter,
		Delta: new(int64),
	}

	// Set up the mock to return the existing metric
	mockGetter.On("Get", "test_metric:counter").Return(existingMetric, true)

	// Test updating the existing metric
	mockSaver.On("Save", "test_metric:counter", existingMetric).Once()

	// Call the method
	updatedMetric := strategy.UpdateMetric(metric)

	// Assert that the delta of the existing metric has been updated
	assert.Equal(t, existingMetric, updatedMetric)
	assert.Equal(t, *existingMetric.Delta, *metric.Delta)

	// Assert the expected method calls
	mockSaver.AssertExpectations(t)
	mockGetter.AssertExpectations(t)
}
