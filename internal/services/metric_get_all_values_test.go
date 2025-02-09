package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetAllRepo is a mock implementation of the GetAllRepo interface
type MockGetAllRepo struct {
	mock.Mock
}

func (m *MockGetAllRepo) GetAll(ctx context.Context) [][]string {
	args := m.Called(ctx)
	return args.Get(0).([][]string)
}

// Test for successful retrieval of all metric values
func TestGetAllMetricValuesService_GetAllMetricValues_Success(t *testing.T) {
	mockRepo := new(MockGetAllRepo)
	service := GetAllMetricValuesService{MetricRepository: mockRepo}

	// Mock the repository to return a list of metrics
	metricsList := [][3]string{
		{"counter", "metric1", "10"},
		{"gauge", "metric2", "20.5"},
	}

	mockRepo.On("GetAll", context.Background()).Return(metricsList)

	// Call the method with context
	metrics := service.GetAllMetricValues(context.Background())

	// Assert no error occurred
	assert.NotNil(t, metrics)
	assert.Len(t, metrics, 2) // Ensure two metrics are returned

	// Check the first metric
	assert.Equal(t, "counter", metrics[0].Type)
	assert.Equal(t, "metric1", metrics[0].Name)
	assert.Equal(t, "10", metrics[0].Value)

	// Check the second metric
	assert.Equal(t, "gauge", metrics[1].Type)
	assert.Equal(t, "metric2", metrics[1].Name)
	assert.Equal(t, "20.5", metrics[1].Value)

	// Assert mock expectations
	mockRepo.AssertExpectations(t)
}

// Test for empty result when no metrics are available
func TestGetAllMetricValuesService_GetAllMetricValues_Empty(t *testing.T) {
	mockRepo := new(MockGetAllRepo)
	service := GetAllMetricValuesService{MetricRepository: mockRepo}

	// Mock the repository to return an empty list
	mockRepo.On("GetAll", context.Background()).Return([][3]string{})

	// Call the method with context
	metrics := service.GetAllMetricValues(context.Background())

	// Assert no error occurred
	assert.NotNil(t, metrics)
	assert.Len(t, metrics, 0) // Ensure the result is empty

	// Assert mock expectations
	mockRepo.AssertExpectations(t)
}
