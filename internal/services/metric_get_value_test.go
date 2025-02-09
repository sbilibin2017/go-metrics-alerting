package services

import (
	"context"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for invalid metric type
func TestGetMetricValueService_GetMetricValue_InvalidType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := GetMetricValueService{MetricRepository: mockRepo}

	req := &types.GetMetricValueRequest{
		Type: types.EmptyString, // Empty type
		Name: "metric1",
	}

	// Call the method with context
	_, err := service.GetMetricValue(context.Background(), req)

	// Assert that an error is returned with the correct code and message
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, errors.ErrInvalidMetricType.Error(), apiErr.Message)
}

// Test for invalid metric name
func TestGetMetricValueService_GetMetricValue_InvalidName(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := GetMetricValueService{MetricRepository: mockRepo}

	req := &types.GetMetricValueRequest{
		Type: "gauge",           // Valid type
		Name: types.EmptyString, // Empty name
	}

	// Call the method with context
	_, err := service.GetMetricValue(context.Background(), req)

	// Assert that an error is returned with the correct code and message
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, errors.ErrInvalidMetricName.Error(), apiErr.Message)
}

// Test for metric not found
func TestGetMetricValueService_GetMetricValue_NotFound(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := GetMetricValueService{MetricRepository: mockRepo}

	req := &types.GetMetricValueRequest{
		Type: "gauge",   // Valid type
		Name: "metric1", // Valid name
	}

	// Mock Get method to return error for not found value
	mockRepo.On("Get", context.Background(), "gauge", "metric1").Return("", errors.ErrValueNotFound)

	// Call the method with context
	result, err := service.GetMetricValue(context.Background(), req)

	// Assert that an error is returned with the correct code and message
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
	assert.Equal(t, errors.ErrValueNotFound.Error(), apiErr.Message)

	// Assert that the result is empty
	assert.Equal(t, types.EmptyString, result)

	// Assert mock expectations
	mockRepo.AssertExpectations(t)
}

// Test for successful metric value retrieval
func TestGetMetricValueService_GetMetricValue_Success(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := GetMetricValueService{MetricRepository: mockRepo}

	req := &types.GetMetricValueRequest{
		Type: "gauge",   // Valid type
		Name: "metric1", // Valid name
	}

	// Mock Get method to return a valid value
	mockRepo.On("Get", context.Background(), "gauge", "metric1").Return("10", nil)

	// Call the method with context
	result, err := service.GetMetricValue(context.Background(), req)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the result is correct
	assert.Equal(t, "10", result)

	// Assert mock expectations
	mockRepo.AssertExpectations(t)
}
