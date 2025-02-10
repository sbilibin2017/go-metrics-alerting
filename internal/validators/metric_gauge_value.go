package validators

import (
	"go-metrics-alerting/internal/errors"
	"strconv"
)

type MetricGaugeValidator struct{}

// Validate проверяет, что значение является корректным числом для Gauge.
func (g *MetricGaugeValidator) Validate(value string) error {
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return errors.ErrInvalidGaugeValue
	}
	return nil
}
