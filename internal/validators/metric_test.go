package validators

import (
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"Valid ID", "metric1", nil},
		{"Empty ID", EmptyString, ErrIDEmpty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateMType(t *testing.T) {
	tests := []struct {
		name    string
		mType   types.MType
		wantErr error
	}{
		{"Valid Counter", types.Counter, nil},
		{"Valid Gauge", types.Gauge, nil},
		{"Empty MType", "", ErrMTypeEmpty},
		{"Invalid MType", "invalid", ErrMTypeInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMType(tt.mType)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateDelta(t *testing.T) {
	var validDelta int64 = 100

	tests := []struct {
		name    string
		mType   types.MType
		delta   *int64
		wantErr error
	}{
		{"Valid Counter with Delta", types.Counter, &validDelta, nil},
		{"Counter without Delta", types.Counter, nil, ErrDeltaEmpty},
		{"Gauge without Delta", types.Gauge, nil, nil}, // Delta не требуется для Gauge
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDelta(tt.mType, tt.delta)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateValue(t *testing.T) {
	var validValue float64 = 99.99

	tests := []struct {
		name    string
		mType   types.MType
		value   *float64
		wantErr error
	}{
		{"Valid Gauge with Value", types.Gauge, &validValue, nil},
		{"Gauge without Value", types.Gauge, nil, ErrValueEmpty},
		{"Counter without Value", types.Counter, nil, nil}, // Value не требуется для Counter
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.mType, tt.value)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
