package validators

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"testing"
)

func TestValidateNonEmptyString(t *testing.T) {
	tests := []struct {
		value    string
		field    string
		expected *types.APIErrorResponse
	}{
		{
			value:    "",
			field:    "metric name",
			expected: &types.APIErrorResponse{Code: http.StatusNotFound, Message: "metric name is required"},
		},
		{
			value:    "metric1",
			field:    "metric name",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			result := ValidateNonEmptyString(tt.value, tt.field)
			if !equalErrors(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateMetricType(t *testing.T) {
	tests := []struct {
		metricType types.MetricType
		expected   *types.APIErrorResponse
	}{
		{
			metricType: types.Counter,
			expected:   nil,
		},
		{
			metricType: types.Gauge,
			expected:   nil,
		},
		{
			metricType: "",
			expected:   &types.APIErrorResponse{Code: http.StatusNotFound, Message: "metric type is required"},
		},
		{
			metricType: "invalid",
			expected:   &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "invalid metric type"},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.metricType), func(t *testing.T) {
			result := ValidateMetricType(tt.metricType)
			if !equalErrors(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateUpdateMetricRequest(t *testing.T) {
	tests := []struct {
		req      *types.UpdateMetricValueRequest
		expected *types.APIErrorResponse
	}{
		{
			req:      &types.UpdateMetricValueRequest{Name: "", Type: types.Counter},
			expected: &types.APIErrorResponse{Code: http.StatusNotFound, Message: "metric name is required"},
		},
		{
			req:      &types.UpdateMetricValueRequest{Name: "metric1", Type: "invalid"},
			expected: &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "invalid metric type"},
		},
		{
			req:      &types.UpdateMetricValueRequest{Name: "metric1", Type: types.Counter},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.req.Name, func(t *testing.T) {
			result := ValidateUpdateMetricRequest(tt.req)
			if !equalErrors(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateCounterValue(t *testing.T) {
	tests := []struct {
		value    string
		expected *types.APIErrorResponse
	}{
		{
			value:    "123",
			expected: nil,
		},
		{
			value:    "invalid",
			expected: &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "invalid counter value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			result := ValidateCounterValue(tt.value)
			if !equalErrors(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateGaugeValue(t *testing.T) {
	tests := []struct {
		value    string
		expected *types.APIErrorResponse
	}{
		{
			value:    "123.45",
			expected: nil,
		},
		{
			value:    "invalid",
			expected: &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "invalid gauge value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			result := ValidateGaugeValue(tt.value)
			if !equalErrors(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateMetricValue(t *testing.T) {
	tests := []struct {
		value      string
		metricType types.MetricType
		expected   *types.APIErrorResponse
	}{
		{
			value:      "123",
			metricType: types.Counter,
			expected:   nil,
		},
		{
			value:      "123.45",
			metricType: types.Gauge,
			expected:   nil,
		},
		{
			value:      "invalid",
			metricType: types.Counter,
			expected:   &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "invalid counter value"},
		},
		{
			value:      "invalid",
			metricType: types.Gauge,
			expected:   &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "invalid gauge value"},
		},
		{
			value:      "123",
			metricType: "invalid",
			expected:   &types.APIErrorResponse{Code: http.StatusBadRequest, Message: "unsupported metric type"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.value+"-"+string(tt.metricType), func(t *testing.T) {
			result := ValidateMetricValue(tt.value, tt.metricType)
			if !equalErrors(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func equalErrors(err1, err2 *types.APIErrorResponse) bool {
	if err1 == nil && err2 == nil {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}
	return err1.Code == err2.Code && err1.Message == err2.Message
}
