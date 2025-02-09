package services

import (
	"context"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricRepository is a mock implementation of the MetricRepository interface
type MockMetricRepository struct {
	mock.Mock
}

func (m *MockMetricRepository) Save(ctx context.Context, metricType, metricName, value string) error {
	args := m.Called(ctx, metricType, metricName, value)
	return args.Error(0)
}

func (m *MockMetricRepository) Get(ctx context.Context, metricType, metricName string) (string, error) {
	args := m.Called(ctx, metricType, metricName)
	return args.String(0), args.Error(1)
}

// Test for invalid value conversion (Value cannot be converted to int64)
func TestUpdateMetricValueService_UpdateMetricValue_InvalidValueConversion(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Counter),
		Name:  "metric1",
		Value: "invalid_value", // Invalid value for conversion
	}

	// Mock getting current value of metric (it doesn't matter in this case)
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("10", nil)

	// Call the method with context
	err := service.UpdateMetricValue(context.Background(), req)

	// Assert that the error occurred
	assert.Error(t, err)

	// Assert the error is of type APIError and has the expected message
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "invalid metric value", apiErr.Message)

	// Ensure that expectations are met
	mockRepo.AssertExpectations(t)
}

// Test for valid Counter type update
func TestUpdateMetricValueService_UpdateMetricValue_ValidCounter(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Counter),
		Name:  "metric1",
		Value: "5",
	}

	// Mock getting current value of metric
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("10", nil)
	// Mock saving updated value
	mockRepo.On("Save", mock.Anything, string(types.Counter), "metric1", "15").Return(nil)

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test for valid Gauge type update
func TestUpdateMetricValueService_UpdateMetricValue_ValidGauge(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Gauge),
		Name:  "metric1",
		Value: "5.5",
	}

	// Mock getting current value of metric
	mockRepo.On("Get", mock.Anything, string(types.Gauge), "metric1").Return("0", nil)
	// Mock saving updated value
	mockRepo.On("Save", mock.Anything, string(types.Gauge), "metric1", "5.5").Return(nil)

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test for missing metric type
func TestUpdateMetricValueService_UpdateMetricValue_MissingType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  types.EmptyString,
		Name:  "metric1",
		Value: "5",
	}

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "metric type is required", apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Test for unsupported metric type
func TestUpdateMetricValueService_UpdateMetricValue_UnsupportedType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  "unsupported_type", // Неподдерживаемый тип
		Name:  "metric1",
		Value: "5",
	}

	// Мокирование вызова Get для неподдерживаемого типа
	mockRepo.On("Get", mock.Anything, "unsupported_type", "metric1").Return("", errors.ErrUnsupportedMetricType)

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "unsupported metric type", apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Test for missing metric name
func TestUpdateMetricValueService_UpdateMetricValue_MissingName(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Counter),
		Name:  types.EmptyString,
		Value: "5",
	}

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
	assert.Equal(t, "metric name is required", apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Test for missing metric value
func TestUpdateMetricValueService_UpdateMetricValue_MissingValue(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Counter),
		Name:  "metric1",
		Value: types.EmptyString,
	}

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "metric value is required", apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Test for invalid Counter conversion (invalid current value)
func TestUpdateMetricValueService_UpdateMetricValue_InvalidCounterConversion(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Counter),
		Name:  "metric1",
		Value: "5",
	}

	// Mock invalid current value (cannot convert to int64)
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("invalid", nil)

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "invalid metric value", apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Test for invalid Gauge conversion
func TestUpdateMetricValueService_UpdateMetricValue_InvalidGaugeConversion(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Gauge),
		Name:  "metric1",
		Value: "invalid_value",
	}

	// Mock getting current value of metric
	mockRepo.On("Get", mock.Anything, string(types.Gauge), "metric1").Return("0", nil)

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "invalid metric value", apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Test for saving error
func TestUpdateMetricValueService_UpdateMetricValue_SaveError(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  string(types.Counter),
		Name:  "metric1",
		Value: "5",
	}

	// Mock getting current value of metric
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("10", nil)
	// Simulate error while saving
	mockRepo.On("Save", mock.Anything, string(types.Counter), "metric1", "15").Return(errors.ErrSaveFailed)

	// Use context.Background() as the context for the test
	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
	assert.Equal(t, "metric value is not saved", apiErr.Message)
	mockRepo.AssertExpectations(t)
}
