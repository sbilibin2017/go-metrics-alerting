package validators

import (
	"errors"
	"go-metrics-alerting/internal/apis/converters"
	"go-metrics-alerting/internal/types"
	"regexp"
)

func ValidateMetricID(id string) error {
	if !regexp.MustCompile(`^[a-z]+(?:_[a-z]+)*$`).MatchString(id) {
		return errors.New("invalid metric id")
	}
	return nil
}

func ValidateMetricType(metricType string) error {
	if metricType != string(types.Counter) && metricType != string(types.Gauge) {
		return errors.New("invalid metric type")
	}
	return nil
}

func ValidateCounterMetricString(value string) error {
	_, err := converters.ConvertStringToInt64(value)
	if err != nil {
		return errors.New("invalid counter metric string value")
	}
	return nil
}

func ValidateGaugeMetricString(value string) error {
	_, err := converters.ConvertStringToFloat64(value)
	if err != nil {
		return errors.New("invalid gauge metric string value")
	}
	return nil
}

func ValidateCounterMetric(value *int64) error {
	if value == nil {
		return errors.New("invalid counter metric value")
	}
	return nil
}

func ValidateGaugeMetric(value *float64) error {
	if value == nil {
		return errors.New("invalid gauge metric value")
	}
	return nil
}
