package updatevalue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateGaugeValueStrategyEngine_Update(t *testing.T) {
	engine := &UpdateGaugeValueStrategyEngine{}

	tests := []struct {
		name         string
		newValue     string
		expected     string
		expectingErr bool
	}{
		{
			name:         "valid float value",
			newValue:     "12.34",
			expected:     "12.34",
			expectingErr: false,
		},
		{
			name:         "valid integer value",
			newValue:     "42",
			expected:     "42",
			expectingErr: false,
		},
		{
			name:         "zero value",
			newValue:     "0",
			expected:     "0",
			expectingErr: false,
		},
		{
			name:         "negative float value",
			newValue:     "-8.76",
			expected:     "-8.76",
			expectingErr: false,
		},
		{
			name:         "negative integer value",
			newValue:     "-100",
			expected:     "-100",
			expectingErr: false,
		},
		{
			name:         "small decimal value",
			newValue:     "0.0001",
			expected:     "0.0001",
			expectingErr: false,
		},
		{
			name:         "large float value",
			newValue:     "9999999999.99",
			expected:     "9999999999.99",
			expectingErr: false,
		},
		{
			name:         "invalid string input",
			newValue:     "abc",
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "empty string",
			newValue:     "",
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "whitespace input",
			newValue:     "   ",
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "mixed alphanumeric input",
			newValue:     "123abc",
			expected:     "",
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Update("", tt.newValue)

			if tt.expectingErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, result, "unexpected result")
			}
		})
	}
}
