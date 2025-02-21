package types

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricBodyRequest_ToDomain(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricBodyRequest
		want    *domain.Metrics
	}{
		{
			name: "valid gauge metric",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "gauge",
				Value: ptrFloat64(10.5),
			},
			want: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Gauge,
				Value: ptrFloat64(10.5),
			},
		},
		{
			name: "valid counter metric",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "counter",
				Delta: ptrInt64(5),
			},
			want: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Counter,
				Delta: ptrInt64(5),
			},
		},
		{
			name: "missing value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "gauge",
				Value: nil,
			},
			want: &domain.Metrics{
				ID:    "metric3",
				MType: domain.Gauge,
				Value: nil,
			},
		},
		{
			name: "missing delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "counter",
				Delta: nil,
			},
			want: &domain.Metrics{
				ID:    "metric4",
				MType: domain.Counter,
				Delta: nil,
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

func TestUpdateMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricBodyRequest
		wantErr bool
	}{
		{
			name: "valid request with gauge type",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "gauge",
				Value: ptrFloat64(10.5),
			},
			wantErr: false,
		},
		{
			name: "valid request with counter type",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "counter",
				Delta: ptrInt64(5),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			request: UpdateMetricBodyRequest{
				ID:    "",
				MType: "gauge",
				Value: ptrFloat64(10.5),
			},
			wantErr: true,
		},
		{
			name: "invalid metric type",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "invalidType",
				Value: ptrFloat64(10.5),
			},
			wantErr: true,
		},
		{
			name: "missing value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "gauge",
				Value: nil,
			},
			wantErr: true,
		},
		{
			name: "missing delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric5",
				MType: "counter",
				Delta: nil,
			},
			wantErr: true,
		},
		{
			name: "valid gauge metric",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "gauge",
				Value: ptrFloat64(10.5),
			},
			wantErr: false,
		},
		{
			name: "valid counter metric",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "counter",
				Delta: ptrInt64(5),
			},
			wantErr: false,
		},
		{
			name: "missing value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "gauge",
				Value: nil,
			},
			wantErr: true,
		},
		{
			name: "missing delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "counter",
				Delta: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid metric type",
			request: UpdateMetricBodyRequest{
				ID:    "metric5",
				MType: "unknown",
				Value: ptrFloat64(5.5),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMetricPathRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricPathRequest
		expected string
	}{
		{
			name: "valid request with gauge type and valid value",
			request: UpdateMetricPathRequest{
				ID:    "metric1",
				MType: "gauge",
				Value: "10.5", // valid float value for gauge
			},
			expected: "", // no error expected
		},
		{
			name: "valid request with counter type and valid value",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "counter",
				Value: "5", // valid integer value for counter
			},
			expected: "", // no error expected
		},
		{
			name: "invalid value for gauge type (non-numeric)",
			request: UpdateMetricPathRequest{
				ID:    "metric3",
				MType: "gauge",
				Value: "invalidValue", // invalid float value for gauge
			},
			expected: "invalid value for Gauge metric, must be a valid float", // expect specific error
		},
		{
			name: "invalid value for counter type (non-numeric)",
			request: UpdateMetricPathRequest{
				ID:    "metric4",
				MType: "counter",
				Value: "invalidValue", // invalid integer value for counter
			},
			expected: "invalid value for Counter metric, must be a valid integer", // expect specific error
		},
		{
			name: "invalid metric type",
			request: UpdateMetricPathRequest{
				ID:    "metric5",
				MType: "unknown", // invalid metric type
				Value: "10",      // any valid value
			},
			expected: "invalid metric type", // expect specific error
		},
		{
			name: "empty ID",
			request: UpdateMetricPathRequest{
				ID:    "",      // empty ID
				MType: "gauge", // valid metric type
				Value: "10.5",  // valid value for gauge
			},
			expected: "id is required", // expect specific error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check if the error matches the expected one
			}
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
			name: "valid request with gauge type",
			request: UpdateMetricPathRequest{
				ID:    "metric1",
				MType: "gauge",
				Value: "10.5", // valid float value
			},
			expected: &domain.Metrics{
				ID:    "metric1",
				MType: domain.Gauge,
				Value: ptrFloat64(10.5),
			},
		},
		{
			name: "valid request with counter type",
			request: UpdateMetricPathRequest{
				ID:    "metric2",
				MType: "counter",
				Value: "5", // valid int value
			},
			expected: &domain.Metrics{
				ID:    "metric2",
				MType: domain.Counter,
				Delta: ptrInt64(5),
			},
		},
		{
			name: "invalid value for gauge type",
			request: UpdateMetricPathRequest{
				ID:    "metric3",
				MType: "gauge",
				Value: "invalidValue", // invalid value for float
			},
			expected: &domain.Metrics{
				ID:    "metric3",
				MType: domain.Gauge,
				Value: nil, // since parsing failed
			},
		},
		{
			name: "invalid value for counter type",
			request: UpdateMetricPathRequest{
				ID:    "metric4",
				MType: "counter",
				Value: "invalidValue", // invalid value for int
			},
			expected: &domain.Metrics{
				ID:    "metric4",
				MType: domain.Counter,
				Delta: nil, // since parsing failed
			},
		},
		{
			name: "missing value for gauge type",
			request: UpdateMetricPathRequest{
				ID:    "metric5",
				MType: "gauge",
				Value: "", // empty string should be invalid
			},
			expected: &domain.Metrics{
				ID:    "metric5",
				MType: domain.Gauge,
				Value: nil, // since parsing failed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.ToDomain()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  GetMetricBodyRequest
		expected string
	}{
		{
			name: "valid request with valid ID and valid metric type",
			request: GetMetricBodyRequest{
				ID:    "metric1",
				MType: "gauge", // valid metric type
			},
			expected: "", // no error expected
		},
		{
			name: "valid request with valid ID and counter type",
			request: GetMetricBodyRequest{
				ID:    "metric2",
				MType: "counter", // valid metric type
			},
			expected: "", // no error expected
		},
		{
			name: "invalid metric type",
			request: GetMetricBodyRequest{
				ID:    "metric3",
				MType: "unknown", // invalid metric type
			},
			expected: "invalid metric type", // expect specific error
		},
		{
			name: "empty ID",
			request: GetMetricBodyRequest{
				ID:    "",      // empty ID
				MType: "gauge", // valid metric type
			},
			expected: "id is required", // expect specific error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check if the error matches the expected one
			}
		})
	}
}

func TestGetMetricPathRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  GetMetricPathRequest
		expected string
	}{
		{
			name: "valid request with valid ID and valid metric type",
			request: GetMetricPathRequest{
				ID:    "metric1",
				MType: "gauge", // valid metric type
			},
			expected: "", // no error expected
		},
		{
			name: "valid request with valid ID and counter type",
			request: GetMetricPathRequest{
				ID:    "metric2",
				MType: "counter", // valid metric type
			},
			expected: "", // no error expected
		},
		{
			name: "invalid metric type",
			request: GetMetricPathRequest{
				ID:    "metric3",
				MType: "unknown", // invalid metric type
			},
			expected: "invalid metric type", // expect specific error
		},
		{
			name: "empty ID",
			request: GetMetricPathRequest{
				ID:    "",      // empty ID
				MType: "gauge", // valid metric type
			},
			expected: "id is required", // expect specific error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check if the error matches the expected one
			}
		})
	}
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}
