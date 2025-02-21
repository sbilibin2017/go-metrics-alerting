package types

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricBodyRequest
		expected int
	}{
		{
			name: "valid counter request",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
				Delta: ptrInt64(10),
			},
			expected: http.StatusOK,
		},
		{
			name: "valid gauge request",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: ptrFloat64(5.5),
			},
			expected: http.StatusOK,
		},
		{
			name: "missing ID",
			request: UpdateMetricBodyRequest{
				MType: "gauge",
				Value: ptrFloat64(5.5),
			},
			expected: http.StatusNotFound,
		},
		{
			name: "invalid metric type",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "invalid",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "missing delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "counter",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "missing value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "gauge",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "empty ID",
			request: UpdateMetricBodyRequest{
				ID:    "",
				MType: "gauge",
				Value: ptrFloat64(1.1),
			},
			expected: http.StatusNotFound,
		},
		{
			name: "zero delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric5",
				MType: "counter",
				Delta: ptrInt64(0),
			},
			expected: http.StatusOK,
		},
		{
			name: "zero value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric6",
				MType: "gauge",
				Value: ptrFloat64(0.0),
			},
			expected: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != nil {
				assert.Equal(t, tt.expected, err.Status)
			} else {
				assert.Equal(t, tt.expected, http.StatusOK)
			}
		})
	}
}

func TestUpdateMetricPathRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricPathRequest
		expected int
	}{
		{
			name: "valid counter request",
			request: UpdateMetricPathRequest{
				ID:    "metric1",
				MType: "counter",
				Value: "10",
			},
			expected: http.StatusOK,
		},
		{
			name: "valid gauge request",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: "42.42",
			},
			expected: http.StatusOK,
		},
		{
			name: "missing ID",
			request: UpdateMetricPathRequest{
				MType: "gauge",
				Value: "5.5",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "invalid metric type",
			request: UpdateMetricPathRequest{
				ID:    "metric3",
				MType: "invalid",
				Value: "100",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid counter value (non-numeric)",
			request: UpdateMetricPathRequest{
				ID:    "metric4",
				MType: "counter",
				Value: "invalid",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid gauge value (non-numeric)",
			request: UpdateMetricPathRequest{
				ID:    "metric5",
				MType: "gauge",
				Value: "invalid",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != nil {
				assert.Equal(t, tt.expected, err.Status)
			} else {
				assert.Equal(t, tt.expected, http.StatusOK)
			}
		})
	}
}

func TestGetMetricPathRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  GetMetricRequest
		expected int
	}{
		{
			name: "valid gauge request",
			request: GetMetricRequest{
				ID:    "metric1",
				MType: "gauge",
			},
			expected: http.StatusOK,
		},
		{
			name: "valid counter request",
			request: GetMetricRequest{
				ID:    "metric2",
				MType: "counter",
			},
			expected: http.StatusOK,
		},
		{
			name: "missing ID",
			request: GetMetricRequest{
				MType: "gauge",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "invalid metric type",
			request: GetMetricRequest{
				ID:    "metric3",
				MType: "invalid",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "empty ID",
			request: GetMetricRequest{
				ID:    "",
				MType: "gauge",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "empty metric type",
			request: GetMetricRequest{
				ID:    "metric4",
				MType: "",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != nil {
				assert.Equal(t, tt.expected, err.Status)
			} else {
				assert.Equal(t, tt.expected, http.StatusOK)
			}
		})
	}
}

func TestUpdateMetricBodyRequest_ToDomain(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricBodyRequest
		expected *domain.Metrics
	}{
		{
			name: "valid counter conversion",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
				Delta: ptrInt64(10),
			},
			expected: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Counter,
				Delta: ptrInt64(10),
				Value: nil,
			},
		},
		{
			name: "valid gauge conversion",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: ptrFloat64(42.42),
			},
			expected: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Gauge,
				Delta: nil,
				Value: ptrFloat64(42.42),
			},
		},
		{
			name: "valid counter with zero delta",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "counter",
				Delta: ptrInt64(0),
			},
			expected: &domain.Metrics{
				ID:    "metric3",
				MType: domain.Counter,
				Delta: ptrInt64(0),
				Value: nil,
			},
		},
		{
			name: "valid gauge with zero value",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "gauge",
				Value: ptrFloat64(0.0),
			},
			expected: &domain.Metrics{
				ID:    "metric4",
				MType: domain.Gauge,
				Delta: nil,
				Value: ptrFloat64(0.0),
			},
		},
		{
			name: "valid counter with nil delta and value",
			request: UpdateMetricBodyRequest{
				ID:    "metric5",
				MType: "counter",
			},
			expected: &domain.Metrics{
				ID:    "metric5",
				MType: domain.Counter,
				Delta: nil,
				Value: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.ToDomain()
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.MType, result.MType)
			assert.Equal(t, tt.expected.Delta, result.Delta)
			assert.Equal(t, tt.expected.Value, result.Value)
		})
	}
}

func TestUpdateMetricPathRequest_ToDomain(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricPathRequest
		expected *domain.Metrics
	}{
		{
			name: "valid counter conversion",
			request: UpdateMetricPathRequest{
				ID:    "metric1",
				MType: "counter",
				Value: "100",
			},
			expected: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Counter,
				Delta: ptrInt64(100),
				Value: nil,
			},
		},
		{
			name: "valid gauge conversion",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: "42.42",
			},
			expected: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Gauge,
				Delta: nil,
				Value: ptrFloat64(42.42),
			},
		},
		{
			name: "invalid counter value (non-numeric)",
			request: UpdateMetricPathRequest{
				ID:    "metric3",
				MType: "counter",
				Value: "invalid",
			},
			expected: &domain.Metrics{
				ID:    "metric3",
				MType: domain.Counter,
				Delta: nil, // Ошибка парсинга, delta остаётся nil
				Value: nil,
			},
		},
		{
			name: "invalid gauge value (non-numeric)",
			request: UpdateMetricPathRequest{
				ID:    "metric4",
				MType: "gauge",
				Value: "invalid",
			},
			expected: &domain.Metrics{
				ID:    "metric4",
				MType: domain.Gauge,
				Delta: nil,
				Value: nil, // Ошибка парсинга, value остаётся nil
			},
		},
		{
			name: "zero gauge value",
			request: UpdateMetricPathRequest{
				ID:    "metric5",
				MType: "gauge",
				Value: "0.0",
			},
			expected: &domain.Metrics{
				ID:    "metric5",
				MType: domain.Gauge,
				Delta: nil,
				Value: ptrFloat64(0.0),
			},
		},
		{
			name: "zero counter value",
			request: UpdateMetricPathRequest{
				ID:    "metric6",
				MType: "counter",
				Value: "0",
			},
			expected: &domain.Metrics{
				ID:    "metric6",
				MType: domain.Counter,
				Delta: ptrInt64(0),
				Value: nil,
			},
		},
		{
			name: "empty value",
			request: UpdateMetricPathRequest{
				ID:    "metric7",
				MType: "gauge",
				Value: "",
			},
			expected: &domain.Metrics{
				ID:    "metric7",
				MType: domain.Gauge,
				Delta: nil,
				Value: nil, // Пустое значение не может быть преобразовано
			},
		},
		{
			name: "invalid metric type",
			request: UpdateMetricPathRequest{
				ID:    "metric8",
				MType: "unknown",
				Value: "123",
			},
			expected: &domain.Metrics{
				ID:    "metric8",
				MType: "unknown", // Остаётся некорректным, но парсинг не происходит
				Delta: nil,
				Value: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.ToDomain()
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.MType, result.MType)
			assert.Equal(t, tt.expected.Delta, result.Delta)
			assert.Equal(t, tt.expected.Value, result.Value)
		})
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrFloat64(f float64) *float64 {
	return &f
}
