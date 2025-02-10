package validators

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"testing"
)

func TestMetricValueValidator_Validate(t *testing.T) {
	tests := []struct {
		metricValue string
		expected    error
	}{
		{"123.45", nil},                  // валидное значение метрики
		{"0", nil},                       // валидное значение метрики
		{"-123.45", nil},                 // валидное отрицательное значение
		{"", errors.ErrEmptyMetricValue}, // пустое значение метрики
		{types.EmptyString, errors.ErrEmptyMetricValue}, // значение, равное EmptyString
	}

	validator := &MetricValueValidator{}

	for _, test := range tests {
		t.Run(test.metricValue, func(t *testing.T) {
			err := validator.Validate(test.metricValue)
			if err != nil && err != test.expected {
				t.Errorf("Expected error %v but got %v", test.expected, err)
			} else if err == nil && test.expected != nil {
				t.Errorf("Expected error %v but got nil", test.expected)
			}
		})
	}
}
