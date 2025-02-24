package validators

import (
	"fmt"
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test ValidateEmptyString
func TestValidateEmptyString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		field    string
		expected error
	}{
		{
			name:     "Valid non-empty string",
			value:    "metric1",
			field:    "ID",
			expected: nil,
		},
		{
			name:     "Empty string",
			value:    "",
			field:    "ID",
			expected: fmt.Errorf("ID cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ValidateEmptyString(tt.value, tt.field)
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, tt.expected.Error())
			}
		})
	}
}

// Test ValidateMetricType
func TestValidateMetricType(t *testing.T) {
	tests := []struct {
		name     string
		mtype    string
		expected error
	}{
		{
			name:     "Valid metric type - Counter",
			mtype:    string(domain.Counter),
			expected: nil,
		},
		{
			name:     "Valid metric type - Gauge",
			mtype:    string(domain.Gauge),
			expected: nil,
		},
		{
			name:     "Invalid metric type",
			mtype:    "invalid_type",
			expected: fmt.Errorf("invalid metric type: invalid_type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ValidateMetricType(tt.mtype)
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, tt.expected.Error())
			}
		})
	}
}

// Test ValidateInt64Ptr
func TestValidateInt64Ptr(t *testing.T) {
	tests := []struct {
		name     string
		value    *int64
		field    string
		expected error
	}{
		{
			name:     "Valid non-nil int64 pointer",
			value:    new(int64),
			field:    "Delta",
			expected: nil,
		},
		{
			name:     "Nil int64 pointer",
			value:    nil,
			field:    "Delta",
			expected: fmt.Errorf("Delta cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ValidateInt64Ptr(tt.value, tt.field)
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, tt.expected.Error())
			}
		})
	}
}

// Test ValidateFloat64Ptr
func TestValidateFloat64Ptr(t *testing.T) {
	tests := []struct {
		name     string
		value    *float64
		field    string
		expected error
	}{
		{
			name:     "Valid non-nil float64 pointer",
			value:    new(float64),
			field:    "Value",
			expected: nil,
		},
		{
			name:     "Nil float64 pointer",
			value:    nil,
			field:    "Value",
			expected: fmt.Errorf("Value cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ValidateFloat64Ptr(tt.value, tt.field)
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, tt.expected.Error())
			}
		})
	}
}
