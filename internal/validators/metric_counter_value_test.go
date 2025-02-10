package validators

import (
	"go-metrics-alerting/internal/errors"
	"testing"
)

func TestMetricCounterValidator_Validate(t *testing.T) {
	tests := []struct {
		value    string
		expected error
	}{
		{"123", nil},  // валидное число
		{"0", nil},    // валидное число
		{"-123", nil}, // отрицательное число
		{"123abc", errors.ErrInvalidCounterValue}, // невалидная строка
		{"abc", errors.ErrInvalidCounterValue},    // невалидная строка
		{"", errors.ErrInvalidCounterValue},       // пустая строка
	}

	validator := &MetricCounterValidator{}

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
