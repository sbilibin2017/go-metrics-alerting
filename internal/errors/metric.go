package errors

import e "errors"

// Ошибки валидации.
var (
	ErrEmptyMetricName       = e.New("metric name is required")
	ErrEmptyMetricType       = e.New("metric type is required")
	ErrEmptyMetricValue      = e.New("metric value is required")
	ErrInvalidCounterValue   = e.New("invalid counter value")
	ErrInvalidGaugeValue     = e.New("invalid gauge value")
	ErrUnsupportedMetricType = e.New("unsupported metric type")
	ErrMetricNotFound        = e.New("metric not found")
)
