package numberprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64Processor_Parse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expected     int64
		expectingErr bool
	}{
		{
			name:         "valid positive number",
			input:        "12345",
			expected:     12345,
			expectingErr: false,
		},
		{
			name:         "valid negative number",
			input:        "-12345",
			expected:     -12345,
			expectingErr: false,
		},
		{
			name:         "zero value",
			input:        "0",
			expected:     0,
			expectingErr: false,
		},
		{
			name:         "invalid string",
			input:        "abc",
			expected:     0,
			expectingErr: true,
		},
		{
			name:         "empty string",
			input:        "",
			expected:     0,
			expectingErr: true,
		},
	}

	processor := NewInt64ProcessorEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Parse(tt.input)

			if tt.expectingErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, result, "unexpected result")
			}
		})
	}
}

func TestInt64ParserFormatter_Format(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "positive number",
			input:    12345,
			expected: "12345",
		},
		{
			name:     "negative number",
			input:    -12345,
			expected: "-12345",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
	}

	processor := NewInt64ProcessorEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Format(tt.input)
			assert.Equal(t, tt.expected, result, "unexpected result")
		})
	}
}
