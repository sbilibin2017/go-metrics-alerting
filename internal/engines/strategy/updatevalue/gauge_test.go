package updatevalue

import (
	"go-metrics-alerting/internal/engines/numberprocessor"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateGaugeValueStrategyEngine_Update(t *testing.T) {
	// Создаем экземпляр обработчика с Float64Processor для обработки значений типа float64
	engine := NewUpdateGaugeValueStrategyEngine(numberprocessor.Float64ProcessorEngine{})

	tests := []struct {
		name         string
		newValue     string
		expected     string
		expectingErr bool
	}{
		{
			name:         "valid float64 value",
			newValue:     "10.5",
			expected:     "10.5", // Проверка с корректным значением типа float64
			expectingErr: false,
		},
		{
			name:         "valid float64 value with scientific notation",
			newValue:     "1e3", // 1 * 10^3 = 1000
			expected:     "1000",
			expectingErr: false,
		},
		{
			name:         "negative float64 value",
			newValue:     "-10.5",
			expected:     "-10.5", // Проверка с отрицательным значением
			expectingErr: false,
		},
		{
			name:         "invalid float64 value",
			newValue:     "abc", // Невалидное значение
			expected:     "",
			expectingErr: true, // Мы ожидаем ошибку при попытке распарсить "abc"
		},
		{
			name:         "empty value",
			newValue:     "", // Пустое значение
			expected:     "",
			expectingErr: true, // Ожидаем ошибку, так как пустая строка не является корректным числом
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Вызов метода Update
			result, err := engine.Update("", tt.newValue)

			// Проверка на ошибку
			if tt.expectingErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, result, "unexpected result")
			}
		})
	}
}
