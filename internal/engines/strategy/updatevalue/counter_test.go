package updatevalue

import (
	"go-metrics-alerting/internal/engines/numberprocessor"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCounterValueStrategyEngine_Update(t *testing.T) {
	// Создаем экземпляр обработчика с Int64Processor для обработки int64 значений
	engine := NewUpdateCounterValueStrategyEngine(numberprocessor.Int64ProcessorEngine{})

	tests := []struct {
		name         string
		currentValue string
		newValue     string
		expected     string
		expectingErr bool
	}{
		{
			name:         "valid values",
			currentValue: "10",
			newValue:     "5",
			expected:     "15", // 10 + 5 = 15
			expectingErr: false,
		},
		{
			name:         "valid values with negative number",
			currentValue: "-10",
			newValue:     "5",
			expected:     "-5", // -10 + 5 = -5
			expectingErr: false,
		},
		{
			name:         "increment zero",
			currentValue: "0",
			newValue:     "5",
			expected:     "5", // 0 + 5 = 5
			expectingErr: false,
		},
		{
			name:         "invalid current value",
			currentValue: "abc",
			newValue:     "10",
			expected:     "",
			expectingErr: true, // invalid current value should raise an error
		},
		{
			name:         "invalid new value",
			currentValue: "10",
			newValue:     "xyz",
			expected:     "",
			expectingErr: true, // invalid new value should raise an error
		},
		{
			name:         "empty input values",
			currentValue: "",
			newValue:     "",
			expected:     "",
			expectingErr: true, // invalid empty values should raise an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Update(tt.currentValue, tt.newValue)

			if tt.expectingErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, result, "unexpected result")
			}
		})
	}
}
