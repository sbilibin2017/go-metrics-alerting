package validators

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmptyString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "empty string",
			value:   "",
			wantErr: ErrValueCannotBeEmpty,
		},
		{
			name:    "string with spaces",
			value:   "   ",
			wantErr: ErrValueCannotBeEmpty,
		},
		{
			name:    "non-empty string",
			value:   "non-empty",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmptyString(tt.value)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMType(t *testing.T) {
	tests := []struct {
		name    string
		mtype   string
		wantErr error
	}{
		{
			name:    "valid mtype counter",
			mtype:   string(domain.Counter),
			wantErr: nil,
		},
		{
			name:    "valid mtype gauge",
			mtype:   string(domain.Gauge),
			wantErr: nil,
		},
		{
			name:    "invalid mtype",
			mtype:   "invalidType",
			wantErr: ErrInvalidMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMType(tt.mtype)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDelta(t *testing.T) {
	tests := []struct {
		name    string
		mtype   string
		delta   *int64
		wantErr error
	}{
		{
			name:    "counter mtype with delta",
			mtype:   string(domain.Counter),
			delta:   ptrInt64(10),
			wantErr: nil,
		},
		{
			name:    "counter mtype without delta",
			mtype:   string(domain.Counter),
			delta:   nil,
			wantErr: ErrDeltaRequiredForCounter,
		},
		{
			name:    "gauge mtype with delta",
			mtype:   string(domain.Gauge),
			delta:   nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDelta(tt.mtype, tt.delta)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateValue(t *testing.T) {
	tests := []struct {
		name    string
		mtype   string
		value   *float64
		wantErr error
	}{
		{
			name:    "gauge mtype with value",
			mtype:   string(domain.Gauge),
			value:   ptrFloat64(10.5),
			wantErr: nil,
		},
		{
			name:    "gauge mtype without value",
			mtype:   string(domain.Gauge),
			value:   nil,
			wantErr: ErrValueRequiredForGauge,
		},
		{
			name:    "counter mtype with value",
			mtype:   string(domain.Counter),
			value:   nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.mtype, tt.value)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions to create pointers
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrFloat64(f float64) *float64 {
	return &f
}
