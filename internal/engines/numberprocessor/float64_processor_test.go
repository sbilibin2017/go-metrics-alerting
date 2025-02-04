package numberprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64Processor_Parse(t *testing.T) {
	processor := NewFloat64ProcessorEngine()

	tests := []struct {
		name        string
		value       string
		expected    float64
		expectError bool
	}{
		{
			name:        "valid float",
			value:       "123.45",
			expected:    123.45,
			expectError: false,
		},
		{
			name:        "valid integer as float",
			value:       "100",
			expected:    100.0,
			expectError: false,
		},
		{
			name:        "empty string",
			value:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid string",
			value:       "abc",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Parse(tt.value)

			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, result, "unexpected result")
			}
		})
	}
}

func TestFloat64Processor_Format(t *testing.T) {
	processor := NewFloat64ProcessorEngine()

	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{
			name:     "positive float",
			value:    123.45,
			expected: "123.45",
		},
		{
			name:     "integer as float",
			value:    100.0,
			expected: "100",
		},
		{
			name:     "negative float",
			value:    -456.78,
			expected: "-456.78",
		},
		{
			name:     "zero",
			value:    0.0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Format(tt.value)
			assert.Equal(t, tt.expected, result, "unexpected result")
		})
	}
}
