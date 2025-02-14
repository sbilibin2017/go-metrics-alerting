package types

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест для validateName
func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *APIErrorResponse
	}{
		{
			name:     "Valid name",
			input:    "valid_metric_name",
			expected: nil,
		},
		{
			name:  "Empty name",
			input: "",
			expected: &APIErrorResponse{
				Code:    http.StatusNotFound, // Исправлено
				Message: "metric name is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateName(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// Тест для validateType
func TestValidateType(t *testing.T) {
	tests := []struct {
		name     string
		input    MetricType
		expected *APIErrorResponse
	}{
		{
			name:     "Valid type - Counter",
			input:    Counter,
			expected: nil,
		},
		{
			name:     "Valid type - Gauge",
			input:    Gauge,
			expected: nil,
		},
		{
			name:  "Invalid type",
			input: "invalid_type",
			expected: &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid metric type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateType(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// Тест для Validate метода UpdateMetricValueRequest
func TestUpdateMetricValueRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricValueRequest
		expected *APIErrorResponse
	}{
		{
			name: "Valid counter request",
			request: UpdateMetricValueRequest{
				Type:  Counter,
				Name:  "metric1",
				Value: "123",
			},
			expected: nil,
		},
		{
			name: "Empty name",
			request: UpdateMetricValueRequest{
				Type:  Counter,
				Name:  "",
				Value: "123",
			},
			expected: &APIErrorResponse{
				Code:    http.StatusNotFound, // Исправлено
				Message: "metric name is required",
			},
		},
		{
			name: "Invalid type",
			request: UpdateMetricValueRequest{
				Type:  "invalid_type",
				Name:  "metric2",
				Value: "123",
			},
			expected: &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid metric type",
			},
		},
		{
			name: "Invalid counter value",
			request: UpdateMetricValueRequest{
				Type:  Counter,
				Name:  "metric3",
				Value: "abc",
			},
			expected: &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid counter value",
			},
		},
		{
			name: "Invalid gauge value",
			request: UpdateMetricValueRequest{
				Type:  Gauge,
				Name:  "metric4",
				Value: "abc",
			},
			expected: &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid gauge value",
			},
		},
		{
			name: "Valid gauge request",
			request: UpdateMetricValueRequest{
				Type:  Gauge,
				Name:  "metric5",
				Value: "123.45",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.Validate()
			assert.Equal(t, tt.expected, got)
		})
	}
}

// Тест для Validate метода GetMetricValueRequest
func TestGetMetricValueRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  GetMetricValueRequest
		expected *APIErrorResponse
	}{
		{
			name: "Valid request",
			request: GetMetricValueRequest{
				Type: Counter,
				Name: "metric1",
			},
			expected: nil,
		},
		{
			name: "Empty name",
			request: GetMetricValueRequest{
				Type: Counter,
				Name: "",
			},
			expected: &APIErrorResponse{
				Code:    http.StatusNotFound, // Исправлено
				Message: "metric name is required",
			},
		},
		{
			name: "Invalid type",
			request: GetMetricValueRequest{
				Type: "invalid_type",
				Name: "metric2",
			},
			expected: &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid metric type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.Validate()
			assert.Equal(t, tt.expected, got)
		})
	}
}
