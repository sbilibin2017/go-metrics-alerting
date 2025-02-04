package updatevalue

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNumber_Success(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  interface{}
		expectErr bool
	}{
		{
			name:      "parse int64",
			input:     "123",
			expected:  int64(123),
			expectErr: false,
		},
		{
			name:      "parse float64",
			input:     "123.45",
			expected:  float64(123.45),
			expectErr: false,
		},
		{
			name:      "parse invalid int64",
			input:     "abc",
			expected:  int64(0),
			expectErr: true,
		},
		{
			name:      "parse invalid float64",
			input:     "abc.123",
			expected:  float64(0),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			var err error
			if _, ok := tt.expected.(int64); ok {
				result, err = parseNumber[int64](tt.input)
			} else if _, ok := tt.expected.(float64); ok {
				result, err = parseNumber[float64](tt.input)
			}

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseNumber_ZeroValueAndError(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectInt int64
		expectFlt float64
		expectErr error
	}{
		{
			name:      "invalid int input returns zero and ErrSyntax",
			input:     "test",
			expectInt: 0,
			expectFlt: 0.0,
			expectErr: strconv.ErrSyntax,
		},
		{
			name:      "empty string returns zero and ErrSyntax",
			input:     "",
			expectInt: 0,
			expectFlt: 0.0,
			expectErr: strconv.ErrSyntax,
		},
		{
			name:      "whitespace returns zero and ErrSyntax",
			input:     "   ",
			expectInt: 0,
			expectFlt: 0.0,
			expectErr: strconv.ErrSyntax,
		},
		{
			name:      "non-numeric characters return zero and ErrSyntax",
			input:     "abc123",
			expectInt: 0,
			expectFlt: 0.0,
			expectErr: strconv.ErrSyntax,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_int64", func(t *testing.T) {
			result, err := parseNumber[int64](tt.input)
			assert.Equal(t, tt.expectInt, result, "expected zero value for int64")
			assert.ErrorIs(t, err, tt.expectErr, "expected strconv.ErrSyntax for int64")
		})

		t.Run(tt.name+"_float64", func(t *testing.T) {
			result, err := parseNumber[float64](tt.input)
			assert.Equal(t, tt.expectFlt, result, "expected zero value for float64")
			assert.ErrorIs(t, err, tt.expectErr, "expected strconv.ErrSyntax for float64")
		})
	}
}

func TestFormatNumber_Int64_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "format int64",
			input:    int64(123),
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber[int64](tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatNumber_Float64_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "format int64",
			input:    float64(123.1),
			expected: "123.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber[float64](tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
