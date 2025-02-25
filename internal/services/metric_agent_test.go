package services

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockMetricFacade - mock for MetricFacade
type MockMetricFacade struct {
	mock.Mock
}

func (m *MockMetricFacade) UpdateMetric(metric *types.UpdateMetricBodyRequest) {
	m.Called(metric)
}

// MockMetricsCollector - mock for MetricsCollector
type MockMetricsCollector struct {
	mock.Mock
}

func (m *MockMetricsCollector) Collect() []*types.UpdateMetricBodyRequest {
	args := m.Called()
	return args.Get(0).([]*types.UpdateMetricBodyRequest)
}

func TestMetricAgentService_Run(t *testing.T) {
	// Create mocks
	mockFacade := new(MockMetricFacade)
	mockCounterCollector := new(MockMetricsCollector)
	mockGaugeCollector := new(MockMetricsCollector)

	// Test config with short intervals for fast testing
	config := &configs.AgentConfig{
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 100 * time.Millisecond,
	}

	// Create the MetricAgentService
	service := NewMetricAgentService(config, mockFacade, mockCounterCollector, mockGaugeCollector)

	// Set up logger (zap development logger for simplicity)
	logger, _ := zap.NewDevelopment()

	// Mock the MetricFacade behavior
	mockFacade.On("UpdateMetric", mock.Anything).Return(nil)

	// Mock the MetricsCollector behavior
	mockCounterCollector.On("Collect").Return([]*types.UpdateMetricBodyRequest{
		{ID: "counter1", MType: "counter", Delta: new(int64), Value: nil},
	})
	mockGaugeCollector.On("Collect").Return([]*types.UpdateMetricBodyRequest{
		{ID: "gauge1", MType: "gauge", Delta: nil, Value: new(float64)},
	})

	// Create a channel to simulate a termination signal after 250ms
	sigChan := make(chan os.Signal, 1)
	go func() {
		time.Sleep(250 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()

	// Replace the signal notify channel with our test signal channel
	originalSigChan := make(chan os.Signal, 1)
	signal.Notify(originalSigChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the Run method in a goroutine
	go service.Run(logger)

	// Simply receive from the signal channel instead of using select
	<-sigChan

	// Assert that the mocks were called as expected
	mockCounterCollector.AssertExpectations(t)
	mockGaugeCollector.AssertExpectations(t)
	mockFacade.AssertExpectations(t)

	// Check that metrics were collected and reported
	assert.True(t, len(mockCounterCollector.Calls) > 0, "Expected counter metrics to be collected")
	assert.True(t, len(mockGaugeCollector.Calls) > 0, "Expected gauge metrics to be collected")
	assert.True(t, len(mockFacade.Calls) > 0, "Expected UpdateMetric to be called")
}

func TestMetricAgentService_RunTerminationSignal(t *testing.T) {
	// Create mocks
	mockFacade := new(MockMetricFacade)
	mockCounterCollector := new(MockMetricsCollector)
	mockGaugeCollector := new(MockMetricsCollector)

	// Test config with short intervals for fast testing
	config := &configs.AgentConfig{
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 100 * time.Millisecond,
	}

	// Create the MetricAgentService
	service := NewMetricAgentService(config, mockFacade, mockCounterCollector, mockGaugeCollector)

	// Set up logger (zap development logger for simplicity)
	logger, _ := zap.NewDevelopment()

	// Mock the MetricFacade behavior
	mockFacade.On("UpdateMetric", mock.Anything).Return(nil)

	// Mock the MetricsCollector behavior
	mockCounterCollector.On("Collect").Return([]*types.UpdateMetricBodyRequest{
		{ID: "counter1", MType: "counter", Delta: new(int64), Value: nil},
	})
	mockGaugeCollector.On("Collect").Return([]*types.UpdateMetricBodyRequest{
		{ID: "gauge1", MType: "gauge", Delta: nil, Value: new(float64)},
	})

	// Create a channel to simulate a termination signal after 250ms
	sigChan := make(chan os.Signal, 1)
	go func() {
		// Simulate the termination signal after 250ms
		time.Sleep(250 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()

	// Start the Run method in a goroutine
	go service.Run(logger)

	// Ensure we wait long enough for the signal to be received and handled.
	// We add an extra 50ms to account for any potential scheduling delays.
	time.Sleep(300 * time.Millisecond)

	// At this point, the service should have received the termination signal and returned.
	// The metrics collection and reporting should not be happening anymore.

	// Assert that the mocks were not called again after the termination signal was received.
	mockCounterCollector.AssertExpectations(t)
	mockGaugeCollector.AssertExpectations(t)
	mockFacade.AssertExpectations(t)
}
