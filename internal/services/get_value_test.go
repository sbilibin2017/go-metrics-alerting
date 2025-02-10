package services

import (
	"context"
	e "errors"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricStorageGetter - мок для интерфейса MetricStorageGetter
type MockMetricStorageGetter struct {
	mock.Mock
}

func (m *MockMetricStorageGetter) Get(ctx context.Context, metricType string, metricName string) (string, error) {
	args := m.Called(ctx, metricType, metricName)
	return args.String(0), args.Error(1)
}

func TestGetMetricValueService_GetMetricValue_Success(t *testing.T) {
	mockRepo := new(MockMetricStorageGetter)
	service := &GetMetricValueService{MetricRepository: mockRepo}

	req := &types.GetMetricValueRequest{
		Type: "gauge",
		Name: "metric1",
	}

	mockRepo.On("Get", mock.Anything, "gauge", "metric1").Return("42.42", nil)

	value, err := service.GetMetricValue(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "42.42", value)
	mockRepo.AssertExpectations(t)
}

func TestGetMetricValueService_GetMetricValue_NotFound(t *testing.T) {
	mockRepo := new(MockMetricStorageGetter)
	service := &GetMetricValueService{MetricRepository: mockRepo}

	req := &types.GetMetricValueRequest{
		Type: "counter",
		Name: "unknown_metric",
	}

	mockRepo.On("Get", mock.Anything, "counter", "unknown_metric").Return("", e.New("not found"))

	value, err := service.GetMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
	assert.Equal(t, errors.ErrMetricNotFound.Error(), apiErr.Message)
	assert.Empty(t, value)
	mockRepo.AssertExpectations(t)
}
