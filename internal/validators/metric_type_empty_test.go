package validators

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"testing"
)

func TestMetricTypeValidator_Validate(t *testing.T) {
	tests := []struct {
		metricType string
		expected   error
	}{
		{"counter", nil},                               // валидный тип метрики
		{"gauge", nil},                                 // валидный тип метрики
		{"", errors.ErrEmptyMetricType},                // пустой тип метрики
		{types.EmptyString, errors.ErrEmptyMetricType}, // значение, равное EmptyString
	}

	validator := &MetricTypeValidator{}

	for _, test := range tests {
		t.Run(test.metricType, func(t *testing.T) {
			err := validator.Validate(test.metricType)
			if err != nil && err != test.expected {
				t.Errorf("Expected error %v but got %v", test.expected, err)
			} else if err == nil && test.expected != nil {
				t.Errorf("Expected error %v but got nil", test.expected)
			}
		})
	}
}
