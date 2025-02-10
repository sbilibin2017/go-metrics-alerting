package validators

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"testing"
)

func TestMetricNameValidator_Validate(t *testing.T) {
	tests := []struct {
		metricName string
		expected   error
	}{
		{"valid_metric_name", nil},                     // валидное имя метрики
		{"", errors.ErrEmptyMetricName},                // пустое имя метрики
		{types.EmptyString, errors.ErrEmptyMetricName}, // значение, равное EmptyString
	}

	validator := &MetricNameValidator{}

	for _, test := range tests {
		t.Run(test.metricName, func(t *testing.T) {
			err := validator.Validate(test.metricName)
			if err != nil && err != test.expected {
				t.Errorf("Expected error %v but got %v", test.expected, err)
			} else if err == nil && test.expected != nil {
				t.Errorf("Expected error %v but got nil", test.expected)
			}
		})
	}
}
