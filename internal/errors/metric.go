package errors

import e "errors"

var (
	ErrInvalidValue          = e.New("invalid metric value")
	ErrUnsupportedMetricType = e.New("unsupported metric type")
	ErrInvalidMetricType     = e.New("metric type is required")
	ErrInvalidMetricName     = e.New("metric name is required")
	ErrInvalidMetricValue    = e.New("metric value is required")
	ErrSaveFailed            = e.New("metric value is not saved")
)
