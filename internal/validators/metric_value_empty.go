package validators

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
)

// MetricValueValidator реализует ValueValidator для проверки значения метрики.
type MetricValueValidator struct{}

// Validate проверяет, что значение метрики не пустое.
func (v *MetricValueValidator) Validate(metricValue string) error {
	if metricValue == types.EmptyString {
		return errors.ErrEmptyMetricValue
	}
	return nil
}
