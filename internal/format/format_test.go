package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatInt64 тестирует функцию FormatInt64
func TestFormatInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{name: "Positive number", input: 12345, expected: "12345"},
		{name: "Negative number", input: -12345, expected: "-12345"},
		{name: "Zero", input: 0, expected: "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatInt64(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// TestParseInt64 тестирует функцию ParseInt64
func TestParseInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{name: "Valid positive number", input: "12345", expected: 12345, hasError: false},
		{name: "Valid negative number", input: "-12345", expected: -12345, hasError: false},
		{name: "Zero", input: "0", expected: 0, hasError: false},
		{name: "Invalid string", input: "abc", expected: 0, hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt64(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

// TestFormatFloat64 тестирует функцию FormatFloat64
func TestFormatFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{name: "Positive float", input: 123.4567890123456789, expected: "123.45678901234568"},
		{name: "Negative float", input: -123.4567890123456789, expected: "-123.45678901234568"},
		{name: "Zero", input: 0.0, expected: "0"},
		{name: "Small float", input: 0.0000000123456789, expected: "1.23456789e-08"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFloat64(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// TestParseFloat64 тестирует функцию ParseFloat64
func TestParseFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{name: "Valid positive float", input: "123.456", expected: 123.456, hasError: false},
		{name: "Valid negative float", input: "-123.456", expected: -123.456, hasError: false},
		{name: "Zero", input: "0.0", expected: 0.0, hasError: false},
		{name: "Scientific notation", input: "1.234e-05", expected: 1.234e-05, hasError: false},
		{name: "Invalid string", input: "abc", expected: 0.0, hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloat64(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
