package services

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetter is a mock implementation of the Getter interface.
type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) Get(key string) (*domain.Metrics, bool) {
	args := m.Called(key)
	// Return nil for the *domain.Metrics and false if the key does not exist
	if args.Get(0) == nil {
		return nil, false
	}
	return args.Get(0).(*domain.Metrics), args.Bool(1)
}

func TestGetMetricService(t *testing.T) {
	// Create mock getter
	mockGetter := new(MockGetter)

	// Create the GetMetricService with the mock getter
	service := NewGetMetricService(mockGetter)

	// Test case 1: Metric exists in storage
	existingMetric := &domain.Metrics{
		ID:    "counter1",
		MType: domain.Counter,
		Delta: new(int64),
	}

	// Mock the Get method to return the existing metric
	mockGetter.On("Get", "counter1:counter").Return(existingMetric, true)

	// Call GetMetric and assert the result
	metric, err := service.GetMetric("counter1", domain.Counter)

	assert.NoError(t, err)
	assert.Equal(t, existingMetric, metric)
	mockGetter.AssertExpectations(t)

	// Test case 2: Metric does not exist in storage
	mockGetter.On("Get", "unknown:counter").Return(nil, false)

	// Call GetMetric for a non-existent metric
	metric, err = service.GetMetric("unknown", domain.Counter)

	assert.Error(t, err)
	assert.Equal(t, ErrMetricNotFound, err)
	assert.Nil(t, metric)
	mockGetter.AssertExpectations(t)

	// Test case 3: Metric exists with a different type
	existingGauge := &domain.Metrics{
		ID:    "gauge1",
		MType: domain.Gauge,
	}

	// Mock the Get method to return the existing gauge metric
	mockGetter.On("Get", "gauge1:gauge").Return(existingGauge, true)

	// Call GetMetric for a gauge metric
	metric, err = service.GetMetric("gauge1", domain.Gauge)

	assert.NoError(t, err)
	assert.Equal(t, existingGauge, metric)
	mockGetter.AssertExpectations(t)
}
