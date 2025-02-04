package updatevalue

import (
	"testing"

	"go-metrics-alerting/internal/engines/numberprocessor"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCounterValueStrategyEngine_Update(t *testing.T) {
	// Создаем экземпляр engine с Int64Processor для обработки int64 значений
	engine := &UpdateCounterValueStrategyEngine[int64]{processor: numberprocessor.Int64Processor{}}

	tests := []struct {
		name         string
		currentValue string
		newValue     string
		expected     string
		expectingErr bool
	}{
		{
			name:         "increment positive values",
			currentValue: "10",
			newValue:     "5",
			expected:     "15",
			expectingErr: false,
		},
		{
			name:         "increment zero",
			currentValue: "0",
			newValue:     "7",
			expected:     "7",
			expectingErr: false,
		},
		{
			name:         "increment negative number",
			currentValue: "-3",
			newValue:     "2",
			expected:     "-1",
			expectingErr: false,
		},
		{
			name:         "adding two negative numbers",
			currentValue: "-10",
			newValue:     "-5",
			expected:     "-15",
			expectingErr: false,
		},
		{
			name:         "large numbers addition",
			currentValue: "1000000000",
			newValue:     "2000000000",
			expected:     "3000000000",
			expectingErr: false,
		},
		{
			name:         "invalid current value",
			currentValue: "abc",
			newValue:     "10",
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "invalid new value",
			currentValue: "10",
			newValue:     "xyz",
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "both values invalid",
			currentValue: "foo",
			newValue:     "bar",
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "empty input values",
			currentValue: "",
			newValue:     "",
			expected:     "",
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Update(tt.currentValue, tt.newValue)

			if tt.expectingErr {
				assert.Error(t, err, "expected an error but got none")
				assert.Equal(t, ErrUnprocessableValue, err, "unexpected error")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, result, "unexpected result")
			}
		})
	}
}
