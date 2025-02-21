package validators

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid string",
			input:    "metric1",
			expected: "", // no error expected
		},
		{
			name:     "empty string",
			input:    "",
			expected: "id is required", // error expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateString(tt.input)
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check error message
			}
		})
	}
}

func TestValidateMType(t *testing.T) {
	tests := []struct {
		name     string
		mType    domain.MType
		expected string
	}{
		{
			name:     "valid counter type",
			mType:    domain.Counter,
			expected: "", // no error expected
		},
		{
			name:     "valid gauge type",
			mType:    domain.Gauge,
			expected: "", // no error expected
		},
		{
			name:     "invalid metric type",
			mType:    "unknown",             // invalid type
			expected: "invalid metric type", // error expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMType(tt.mType)
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check error message
			}
		})
	}
}

func TestValidateDelta(t *testing.T) {
	tests := []struct {
		name     string
		mType    domain.MType
		delta    *int64
		expected string
	}{
		{
			name:     "valid delta for counter",
			mType:    domain.Counter,
			delta:    ptrInt64(10),
			expected: "", // no error expected
		},
		{
			name:     "missing delta for counter",
			mType:    domain.Counter,
			delta:    nil,
			expected: "delta is required for Counter metric", // error expected
		},
		{
			name:     "no delta for gauge",
			mType:    domain.Gauge,
			delta:    nil,
			expected: "", // no error expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDelta(tt.mType, tt.delta)
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check error message
			}
		})
	}
}

func TestValidateValue(t *testing.T) {
	tests := []struct {
		name     string
		mType    domain.MType
		value    *float64
		expected string
	}{
		{
			name:     "valid value for gauge",
			mType:    domain.Gauge,
			value:    ptrFloat64(10.5),
			expected: "", // no error expected
		},
		{
			name:     "missing value for gauge",
			mType:    domain.Gauge,
			value:    nil,
			expected: "value is required for Gauge metric", // error expected
		},
		{
			name:     "no value for counter",
			mType:    domain.Counter,
			value:    nil,
			expected: "", // no error expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.mType, tt.value)
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check error message
			}
		})
	}
}

func TestValidateValueString(t *testing.T) {
	tests := []struct {
		name     string
		mType    domain.MType
		value    string
		expected string
	}{
		{
			name:     "valid value for gauge",
			mType:    domain.Gauge,
			value:    "10.5", // valid float for gauge
			expected: "",     // no error expected
		},
		{
			name:     "invalid value for gauge",
			mType:    domain.Gauge,
			value:    "invalid",                                               // invalid float for gauge
			expected: "invalid value for Gauge metric, must be a valid float", // error expected
		},
		{
			name:     "valid value for counter",
			mType:    domain.Counter,
			value:    "10", // valid int for counter
			expected: "",   // no error expected
		},
		{
			name:     "invalid value for counter",
			mType:    domain.Counter,
			value:    "invalid",                                                   // invalid int for counter
			expected: "invalid value for Counter metric, must be a valid integer", // error expected
		},
		{
			name:     "unsupported metric type",
			mType:    "unknown", // invalid type
			value:    "10",
			expected: "unsupported metric type", // error expected
		},
		{
			name:     "empty value",
			mType:    domain.Counter,
			value:    "",
			expected: "id is required", // check for ValidateString call
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValueString(tt.mType, tt.value)
			if tt.expected == "" {
				require.NoError(t, err) // no error expected
			} else {
				assert.EqualError(t, err, tt.expected) // check error message
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
