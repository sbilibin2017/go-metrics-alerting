package validators

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
)

// MetricTypeValidator реализует TypeValidator для проверки типа метрики.
type MetricTypeValidator struct{}

// Validate проверяет, что тип метрики не пустой.
func (v *MetricTypeValidator) Validate(metricType string) error {
	if metricType == types.EmptyString {
		return errors.ErrEmptyMetricType
	}
	return nil
}
