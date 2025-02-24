package types

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GetMetricRequest.Validate
func TestGetMetricRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  GetMetricRequest
		expected int
	}{
		{
			name: "Valid request",
			request: GetMetricRequest{
				ID:    "metric1",
				MType: "counter",
			},
			expected: http.StatusOK,
		},
		{
			name: "Invalid ID (Empty)",
			request: GetMetricRequest{
				ID:    "",
				MType: "counter",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid MType (Empty)",
			request: GetMetricRequest{
				ID:    "metric2",
				MType: "",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid MType (Invalid Value)",
			request: GetMetricRequest{
				ID:    "metric2",
				MType: "invalid_type",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Valid MType (Counter)",
			request: GetMetricRequest{
				ID:    "metric3",
				MType: "counter",
			},
			expected: http.StatusOK,
		},
		{
			name: "Valid MType (Gauge)",
			request: GetMetricRequest{
				ID:    "metric4",
				MType: "gauge",
			},
			expected: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.request.Validate()
			if actual == nil {
				assert.Equal(t, http.StatusOK, tt.expected)
			} else {
				assert.Equal(t, tt.expected, actual.Status)
			}
		})
	}
}
