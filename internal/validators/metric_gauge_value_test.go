package validators

import (
	"go-metrics-alerting/internal/errors"
	"testing"
)

func TestMetricGaugeValidator_Validate(t *testing.T) {
	tests := []struct {
		value    string
		expected error
	}{
		{"123.45", nil},                         // валидное число с плавающей запятой
		{"0.0", nil},                            // валидное число с плавающей запятой
		{"-123.45", nil},                        // отрицательное число с плавающей запятой
		{"123", nil},                            // валидное целое число
		{"123.0", nil},                          // целое число, но с плавающей точкой
		{"123abc", errors.ErrInvalidGaugeValue}, // невалидная строка
		{"abc", errors.ErrInvalidGaugeValue},    // невалидная строка
		{"", errors.ErrInvalidGaugeValue},       // пустая строка
	}

	validator := &MetricGaugeValidator{}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			err := validator.Validate(test.value)
			if err != nil && err != test.expected {
				t.Errorf("Expected error %v but got %v", test.expected, err)
			} else if err == nil && test.expected != nil {
				t.Errorf("Expected error %v but got nil", test.expected)
			}
		})
	}
}
