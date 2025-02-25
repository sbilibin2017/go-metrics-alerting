package types

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test UpdateMetricPathRequest.Validate
func TestUpdateMetricPathRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricPathRequest
		expected int
	}{
		{
			name: "Valid request",
			request: UpdateMetricPathRequest{
				ID:    "metric1",
				MType: "counter",
				Value: "100",
			},
			expected: http.StatusOK,
		},
		{
			name: "Invalid ID",
			request: UpdateMetricPathRequest{
				ID:    "",
				MType: "counter",
				Value: "100",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid MType (Empty)",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "",
				Value: "100",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid MType (Invalid Value)",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "invalid_type",
				Value: "100",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Valid MType (Counter)",
			request: UpdateMetricPathRequest{
				ID:    "metric3",
				MType: "counter",
				Value: "100",
			},
			expected: http.StatusOK,
		},
		{
			name: "Valid MType (Gauge)",
			request: UpdateMetricPathRequest{
				ID:    "metric4",
				MType: "gauge",
				Value: "50.5",
			},
			expected: http.StatusOK,
		},
		{
			name: "Empty Value for Gauge",
			request: UpdateMetricPathRequest{
				ID:    "metric5",
				MType: "gauge",
				Value: "",
			},
			expected: http.StatusBadRequest,
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

// Test UpdateMetricPathRequest.ToDomain with valid values
func TestUpdateMetricPathRequest_ToDomain(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricPathRequest
		want    *domain.Metrics
	}{
		{
			name: "Convert to domain model for counter with valid value",
			request: UpdateMetricPathRequest{
				ID:    "metric1",
				MType: "counter",
				Value: "100",
			},
			want: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Counter,
				Delta: pathInt64Ptr(100), // ожидаем корректный Delta для counter
				Value: nil,               // Value должно быть nil для counter
			},
		},
		{
			name: "Convert to domain model for gauge with valid value",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: "25.5",
			},
			want: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Gauge,
				Delta: nil,                  // для gauge Delta должно быть nil
				Value: pathFloat64Ptr(25.5), // ожидаем корректное значение для Value
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToDomain()
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper function to create a pointer for int64
func pathInt64Ptr(i int64) *int64 {
	return &i
}

// Helper function to create a pointer for float64
func pathFloat64Ptr(f float64) *float64 {
	return &f
}
