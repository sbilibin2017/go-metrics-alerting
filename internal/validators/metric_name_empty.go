package validators

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
)

// MetricNameValidator реализует NameValidator для проверки имени метрики.
type MetricNameValidator struct{}

// Validate проверяет, что имя метрики не пустое.
func (v *MetricNameValidator) Validate(metricName string) error {
	if metricName == types.EmptyString {
		return errors.ErrEmptyMetricName
	}
	return nil
}
