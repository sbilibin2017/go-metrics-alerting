package validators

import (
	"go-metrics-alerting/internal/errors"
	"strconv"
)

type MetricCounterValidator struct{}

func (c *MetricCounterValidator) Validate(value string) error {
	if _, err := strconv.ParseInt(value, 10, 64); err != nil {
		return errors.ErrInvalidCounterValue
	}
	return nil
}
