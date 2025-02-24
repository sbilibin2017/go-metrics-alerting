package types

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GetAllMetricsResponse.FromDomain
func TestGetAllMetricsResponse_FromDomain(t *testing.T) {
	tests := []struct {
		name     string
		metric   *domain.Metrics
		expected *GetAllMetricsResponse
	}{
		{
			name: "Counter metric",
			metric: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Counter,
				Delta: new(int64),
			},
			expected: &GetAllMetricsResponse{
				ID:    "metric1",
				Value: "0",
			},
		},
		{
			name: "Gauge metric",
			metric: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Gauge,
				Value: new(float64),
			},
			expected: &GetAllMetricsResponse{
				ID:    "metric2",
				Value: "0.000000",
			},
		},
		{
			name: "Counter metric with non-zero value",
			metric: &domain.Metrics{
				ID:    "metric3",
				MType: domain.Counter,
				Delta: func() *int64 { v := int64(123); return &v }(),
			},
			expected: &GetAllMetricsResponse{
				ID:    "metric3",
				Value: "123",
			},
		},
		{
			name: "Gauge metric with non-zero value",
			metric: &domain.Metrics{
				ID:    "metric4",
				MType: domain.Gauge,
				Value: func() *float64 { v := float64(45.67); return &v }(),
			},
			expected: &GetAllMetricsResponse{
				ID:    "metric4",
				Value: "45.670000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &GetAllMetricsResponse{}
			actual := resp.FromDomain(tt.metric)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
