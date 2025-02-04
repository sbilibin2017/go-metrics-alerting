package updatevalue

import (
	"testing"

	"go-metrics-alerting/internal/engines/numberprocessor"

	"github.com/stretchr/testify/assert"
)

func TestUpdateGaugeValueStrategyEngine_Update(t *testing.T) {
	// Создаем экземпляр engine с Int64Processor для обработки int64 значений
	engineInt := &UpdateGaugeValueStrategyEngine[int64]{processor: numberprocessor.Int64Processor{}}
	// Создаем экземпляр engine с Float64Processor для обработки float64 значений
	engineFloat := &UpdateGaugeValueStrategyEngine[float64]{processor: numberprocessor.Float64Processor{}}

	tests := []struct {
		name         string
		engine       interface{}
		newValue     string
		expected     string
		expectingErr bool
	}{
		// Тесты для int64
		{
			name:         "valid int64 value",
			engine:       engineInt,
			newValue:     "10",
			expected:     "10",
			expectingErr: false,
		},
		{
			name:         "invalid int64 value",
			engine:       engineInt,
			newValue:     "abc",
			expected:     "",
			expectingErr: true,
		},

		// Тесты для float64
		{
			name:         "valid float64 value",
			engine:       engineFloat,
			newValue:     "10.5",
			expected:     "10.5",
			expectingErr: false,
		},
		{
			name:         "invalid float64 value",
			engine:       engineFloat,
			newValue:     "xyz",
			expected:     "",
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			var err error

			// Проверка типа engine
			switch e := tt.engine.(type) {
			case *UpdateGaugeValueStrategyEngine[int64]:
				result, err = e.Update("", tt.newValue)
			case *UpdateGaugeValueStrategyEngine[float64]:
				result, err = e.Update("", tt.newValue)
			default:
				t.Fatalf("unexpected engine type %T", e)
			}

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
