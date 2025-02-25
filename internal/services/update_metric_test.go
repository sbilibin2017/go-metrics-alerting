package services

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for UpdateMetricStrategy
type MockUpdateMetricStrategy struct {
	mock.Mock
}

func (m *MockUpdateMetricStrategy) UpdateMetric(metric *domain.Metrics) *domain.Metrics {
	args := m.Called(metric)
	return args.Get(0).(*domain.Metrics)
}

func TestUpdateMetricService(t *testing.T) {
	// Create mock strategies
	mockCounterStrategy := new(MockUpdateMetricStrategy)
	mockGaugeStrategy := new(MockUpdateMetricStrategy)

	// Create the UpdateMetricService
	service := NewUpdateMetricService(mockCounterStrategy, mockGaugeStrategy)

	// Test case 1: Valid Counter Metric
	counterMetric := &domain.Metrics{
		ID:    "counter1",
		MType: domain.Counter,
		Delta: new(int64),
	}

	mockCounterStrategy.On("UpdateMetric", counterMetric).Return(counterMetric)

	// Call UpdateMetric for Counter
	updatedMetric, err := service.UpdateMetric(counterMetric)

	// Assertions for Counter Metric
	assert.NoError(t, err)
	assert.Equal(t, counterMetric, updatedMetric)
	mockCounterStrategy.AssertExpectations(t)

	// Test case 2: Valid Gauge Metric
	gaugeMetric := &domain.Metrics{
		ID:    "gauge1",
		MType: domain.Gauge,
	}

	mockGaugeStrategy.On("UpdateMetric", gaugeMetric).Return(gaugeMetric)

	// Call UpdateMetric for Gauge
	updatedMetric, err = service.UpdateMetric(gaugeMetric)

	// Assertions for Gauge Metric
	assert.NoError(t, err)
	assert.Equal(t, gaugeMetric, updatedMetric)
	mockGaugeStrategy.AssertExpectations(t)

	// Test case 3: Invalid Metric Type
	unknownMetric := &domain.Metrics{
		ID:    "unknown1",
		MType: "unknown", // Invalid type
	}

	// Call UpdateMetric for an invalid metric type
	updatedMetric, err = service.UpdateMetric(unknownMetric)

	// Assertions for Invalid Metric Type
	assert.Error(t, err)
	assert.Equal(t, ErrUnknownMetricType, err)
	assert.Nil(t, updatedMetric)
}
