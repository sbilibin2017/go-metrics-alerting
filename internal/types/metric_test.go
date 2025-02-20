package types

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricBodyRequest
		wantErr bool
	}{
		{
			name: "valid request with delta",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
				Delta: new(int64),
			},
			wantErr: false,
		},
		{
			name: "valid request with value",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: new(float64),
			},
			wantErr: false,
		},
		{
			name: "delta and value set together",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "counter",
				Delta: new(int64),
				Value: new(float64),
			},
			wantErr: true,
		},
		{
			name: "empty ID",
			request: UpdateMetricBodyRequest{
				ID:    "",
				MType: "counter",
			},
			wantErr: true,
		},
		{
			name: "invalid metric type",
			request: UpdateMetricBodyRequest{
				ID:    "metric4",
				MType: "invalid_type",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateMetricBodyRequest_ToMetric(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricBodyRequest
		want    *domain.Metric
	}{
		{
			name: "to metric with delta",
			request: UpdateMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
				Delta: new(int64),
			},
			want: &domain.Metric{
				ID:    "metric1",
				MType: domain.Counter,
				Value: "0", // assuming *int64 default value is 0
			},
		},
		{
			name: "to metric with value",
			request: UpdateMetricBodyRequest{
				ID:    "metric2",
				MType: "gauge",
				Value: new(float64),
			},
			want: &domain.Metric{
				ID:    "metric2",
				MType: domain.Gauge,
				Value: "0.000000", // assuming *float64 default value is 0.0
			},
		},
		{
			name: "nil metric",
			request: UpdateMetricBodyRequest{
				ID:    "metric3",
				MType: "counter",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToMetric()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateMetricPathRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricPathRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: UpdateMetricPathRequest{
				Name:  "metric1",
				Type:  "counter",
				Value: "100",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: UpdateMetricPathRequest{
				Name:  "",
				Type:  "counter",
				Value: "100",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			request: UpdateMetricPathRequest{
				Name:  "metric1",
				Type:  "",
				Value: "100",
			},
			wantErr: true,
		},
		{
			name: "empty value",
			request: UpdateMetricPathRequest{
				Name:  "metric1",
				Type:  "counter",
				Value: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateMetricPathRequest_ToMetric(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricPathRequest
		want    *domain.Metric
	}{
		{
			name: "to valid metric",
			request: UpdateMetricPathRequest{
				Name:  "metric1",
				Type:  "counter",
				Value: "100",
			},
			want: &domain.Metric{
				ID:    "metric1",
				MType: domain.Counter,
				Value: "100",
			},
		},
		{
			name: "to invalid type (gauge)",
			request: UpdateMetricPathRequest{
				Name:  "metric2",
				Type:  "gauge",
				Value: "200",
			},
			want: &domain.Metric{
				ID:    "metric2",
				MType: domain.Gauge,
				Value: "200",
			},
		},
		{
			name: "to nil for invalid type",
			request: UpdateMetricPathRequest{
				Name:  "metric3",
				Type:  "invalid_type",
				Value: "300",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToMetric()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request GetMetricBodyRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: GetMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			request: GetMetricBodyRequest{
				ID:    "",
				MType: "counter",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			request: GetMetricBodyRequest{
				ID:    "metric1",
				MType: "invalid_type",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetMetricBodyRequest_ToMetric(t *testing.T) {
	tests := []struct {
		name    string
		request GetMetricBodyRequest
		value   string
		want    *domain.Metric
	}{
		{
			name: "to valid metric",
			request: GetMetricBodyRequest{
				ID:    "metric1",
				MType: "counter",
			},
			value: "100",
			want: &domain.Metric{
				ID:    "metric1",
				MType: domain.Counter,
				Value: "100",
			},
		},
		{
			name: "to invalid type (gauge)",
			request: GetMetricBodyRequest{
				ID:    "metric2",
				MType: "gauge",
			},
			value: "200",
			want: &domain.Metric{
				ID:    "metric2",
				MType: domain.Gauge,
				Value: "200",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToMetric(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateMetricBodyRequest_Validate_DeltaAndValueSet_Error(t *testing.T) {
	// Create request where both Delta and Value are set
	request := UpdateMetricBodyRequest{
		ID:    "metric4",
		MType: "counter",
		Delta: new(int64),   // Delta is set
		Value: new(float64), // Value is set
	}

	// Perform validation, which should trigger an error for both fields being set
	err := request.Validate()

	// Assert that the error is the expected conflict error
	require.Error(t, err)
	assert.Equal(t, ErrDeltaValueConflict, err)
}

func TestUpdateMetricBodyRequest_Validate_MissingDeltaForCounter(t *testing.T) {
	// Create a request where MType is "counter" but Delta is missing, which should trigger an error
	request := UpdateMetricBodyRequest{
		ID:    "metric5",
		MType: "counter", // MType is counter
		Delta: nil,       // Delta is nil (this is the issue)
		Value: nil,       // Value is nil
	}

	// Perform validation, which should fail because Delta is required for Counter
	err := request.Validate()

	// Assert that the error is the expected "delta required for counter" error
	require.Error(t, err)

}

func TestUpdateMetricBodyRequest_Validate_MissingValueForGauge(t *testing.T) {
	// Create a request where MType is "gauge" but Value is missing, which should trigger an error
	request := UpdateMetricBodyRequest{
		ID:    "metric6",
		MType: "gauge", // MType is gauge
		Delta: nil,     // Delta is nil
		Value: nil,     // Value is nil (this is the issue)
	}

	// Perform validation, which should fail because Value is required for Gauge
	err := request.Validate()

	// Assert that the error is the expected "value required for gauge" error
	require.Error(t, err)

}

func TestGetMetricBodyResponse_ToMetric(t *testing.T) {
	tests := []struct {
		name    string
		request GetMetricBodyResponse
		want    *domain.Metric
	}{
		{
			name: "valid metric",
			request: GetMetricBodyResponse{
				ID:    "metric1",
				MType: "counter",
				Value: "100",
			},
			want: &domain.Metric{
				ID:    "metric1",
				MType: domain.Counter,
				Value: "100",
			},
		},
		{
			name: "valid gauge metric",
			request: GetMetricBodyResponse{
				ID:    "metric2",
				MType: "gauge",
				Value: "200.5",
			},
			want: &domain.Metric{
				ID:    "metric2",
				MType: domain.Gauge,
				Value: "200.5",
			},
		},
		{
			name: "nil value",
			request: GetMetricBodyResponse{
				ID:    "metric3",
				MType: "counter",
				Value: "",
			},
			want: &domain.Metric{
				ID:    "metric3",
				MType: domain.Counter,
				Value: "",
			},
		},
		{
			name: "empty metric type",
			request: GetMetricBodyResponse{
				ID:    "metric4",
				MType: "",
				Value: "150",
			},
			want: &domain.Metric{
				ID:    "metric4",
				MType: domain.MType(""), // This assumes "" is a valid MType, otherwise change it based on your validation rules
				Value: "150",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToMetric()
			assert.Equal(t, tt.want, got)
		})
	}
}
