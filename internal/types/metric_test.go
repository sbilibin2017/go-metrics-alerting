package types

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/validators"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricBodyRequest
		wantErr error
	}{
		{
			name: "valid request with counter and delta",
			request: UpdateMetricBodyRequest{
				ID:    "metric_1",
				MType: string(domain.Counter),
				Delta: ptrInt64(10),
			},
			wantErr: nil,
		},
		{
			name: "valid request with gauge and value",
			request: UpdateMetricBodyRequest{
				ID:    "metric_2",
				MType: string(domain.Gauge),
				Value: ptrFloat64(100.5),
			},
			wantErr: nil,
		},
		{
			name: "invalid request with empty ID",
			request: UpdateMetricBodyRequest{
				ID:    "",
				MType: string(domain.Counter),
			},
			wantErr: validators.ErrValueCannotBeEmpty,
		},
		{
			name: "invalid request with invalid MType",
			request: UpdateMetricBodyRequest{
				ID:    "metric_3",
				MType: "invalidType",
			},
			wantErr: validators.ErrInvalidMetricType,
		},
		{
			name: "invalid request with missing Delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "metric_4",
				MType: string(domain.Counter),
				Delta: nil,
			},
			wantErr: validators.ErrDeltaRequiredForCounter,
		},
		{
			name: "invalid request with missing Value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "metric_5",
				MType: string(domain.Gauge),
				Value: nil,
			},
			wantErr: validators.ErrValueRequiredForGauge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMetricQueryRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMetricQueryRequest
		wantErr error
	}{
		{
			name: "valid request",
			request: UpdateMetricQueryRequest{
				Name:  "metric_1",
				Type:  string(domain.Counter),
				Value: "100",
			},
			wantErr: nil,
		},
		{
			name: "invalid request with empty Name",
			request: UpdateMetricQueryRequest{
				Name:  "",
				Type:  string(domain.Counter),
				Value: "100",
			},
			wantErr: validators.ErrValueCannotBeEmpty,
		},
		{
			name: "invalid request with empty Type",
			request: UpdateMetricQueryRequest{
				Name:  "metric_2",
				Type:  "",
				Value: "100",
			},
			wantErr: validators.ErrValueCannotBeEmpty,
		},
		{
			name: "invalid request with empty Value",
			request: UpdateMetricQueryRequest{
				Name:  "metric_3",
				Type:  string(domain.Counter),
				Value: "",
			},
			wantErr: validators.ErrValueCannotBeEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMetricBodyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request GetMetricBodyRequest
		wantErr error
	}{
		{
			name: "valid request",
			request: GetMetricBodyRequest{
				ID:    "metric_1",
				MType: string(domain.Counter),
			},
			wantErr: nil,
		},
		{
			name: "invalid request with empty ID",
			request: GetMetricBodyRequest{
				ID:    "",
				MType: string(domain.Counter),
			},
			wantErr: validators.ErrValueCannotBeEmpty,
		},
		{
			name: "invalid request with invalid MType",
			request: GetMetricBodyRequest{
				ID:    "metric_2",
				MType: "invalidType",
			},
			wantErr: validators.ErrInvalidMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrFloat64(f float64) *float64 {
	return &f
}
