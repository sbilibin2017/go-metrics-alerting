package types

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test UpdateMetricBodyRequest.Validate
func TestUpdateMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricBodyRequest
		expected int
	}{
		{
			name: "Valid Counter request with Delta",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
				Delta: new(int64),
			},
			expected: http.StatusOK,
		},
		{
			name: "Invalid ID",
			request: UpdateMetricBodyRequest{
				ID:    "",
				MType: "counter",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid MType",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "invalid_type",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid Delta for Counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "counter",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid Value for Gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "gauge",
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

// Test UpdateMetricBodyRequest.ToDomain
func TestUpdateMetricBodyRequest_ToDomain(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricBodyRequest
		want    *domain.Metrics
	}{
		{
			name: "Convert to domain model for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
				Delta: new(int64),
			},
			want: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Counter,
				Delta: new(int64),
				Value: nil,
			},
		},
		{
			name: "Convert to domain model for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: new(float64),
			},
			want: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Gauge,
				Delta: nil,
				Value: new(float64),
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
