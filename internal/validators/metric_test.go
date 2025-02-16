package validators

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetrics(t *testing.T) {
	tests := []struct {
		name     string
		input    *types.MetricsRequest
		expected *types.APIErrorResponse
	}{
		{
			name: "Valid counter metric with delta",
			input: &types.MetricsRequest{
				ID:    "metric1",
				MType: types.Counter,
				Delta: new(int64),
				Value: nil,
			},
			expected: nil,
		},
		{
			name: "Valid gauge metric with value",
			input: &types.MetricsRequest{
				ID:    "metric2",
				MType: types.Gauge,
				Delta: nil,
				Value: new(float64),
			},
			expected: nil,
		},
		{
			name: "Missing ID",
			input: &types.MetricsRequest{
				ID:    "",
				MType: types.Counter,
				Delta: new(int64),
				Value: nil,
			},
			expected: types.NewAPIErrorResponse(http.StatusNotFound, ErrIDRequired),
		},
		{
			name: "Missing Metric Type",
			input: &types.MetricsRequest{
				ID:    "metric3",
				MType: "",
				Delta: new(int64),
				Value: nil,
			},
			expected: types.NewAPIErrorResponse(http.StatusNotFound, ErrMTypeRequired),
		},
		{
			name: "Invalid Metric Type",
			input: &types.MetricsRequest{
				ID:    "metric4",
				MType: types.MType("invalid"),
				Delta: new(int64),
				Value: nil,
			},
			expected: types.NewAPIErrorResponse(http.StatusBadRequest, ErrInvalidMetricType),
		},
		{
			name: "Missing delta for counter",
			input: &types.MetricsRequest{
				ID:    "metric5",
				MType: types.Counter,
				Delta: nil,
				Value: nil,
			},
			expected: types.NewAPIErrorResponse(http.StatusBadRequest, ErrDeltaRequired),
		},
		{
			name: "Missing value for gauge",
			input: &types.MetricsRequest{
				ID:    "metric6",
				MType: types.Gauge,
				Delta: nil,
				Value: nil,
			},
			expected: types.NewAPIErrorResponse(http.StatusBadRequest, ErrValueRequired),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMetrics(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateMetricValueRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    *types.MetricValueRequest
		expected *types.APIErrorResponse
	}{
		{
			name: "Valid metric value request",
			input: &types.MetricValueRequest{
				ID:    "metric1",
				MType: types.Counter,
			},
			expected: nil,
		},
		{
			name: "Missing ID",
			input: &types.MetricValueRequest{
				ID:    "",
				MType: types.Counter,
			},
			expected: types.NewAPIErrorResponse(http.StatusNotFound, ErrIDRequired),
		},
		{
			name: "Missing Metric Type",
			input: &types.MetricValueRequest{
				ID:    "metric2",
				MType: "",
			},
			expected: types.NewAPIErrorResponse(http.StatusNotFound, ErrMTypeRequired),
		},
		{
			name: "Invalid Metric Type",
			input: &types.MetricValueRequest{
				ID:    "metric3",
				MType: types.MType("invalid"),
			},
			expected: types.NewAPIErrorResponse(http.StatusBadRequest, ErrInvalidMetricType),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMetricValueRequest(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
