package validators

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"strconv"
)

const EmptyString string = ""

// validateNonEmptyString проверяет, что строка не пустая.
func ValidateNonEmptyString(value string, fieldName string) *types.APIErrorResponse {
	if value == EmptyString {
		return &types.APIErrorResponse{
			Code:    http.StatusNotFound,
			Message: fieldName + " is required",
		}
	}
	return nil
}

// validateMetricType проверяет, что тип метрики валиден.
func ValidateMetricType(metricType types.MetricType) *types.APIErrorResponse {
	if err := ValidateNonEmptyString(string(metricType), "metric type"); err != nil {
		return err
	}
	if metricType != types.Counter && metricType != types.Gauge {
		return &types.APIErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid metric type",
		}
	}
	return nil
}

// ValidateUpdateMetricRequest проверяет запрос на обновление метрики.
func ValidateUpdateMetricRequest(req *types.UpdateMetricValueRequest) *types.APIErrorResponse {
	if err := ValidateNonEmptyString(req.Name, "metric name"); err != nil {
		return err
	}
	if err := ValidateMetricType(req.Type); err != nil {
		return err
	}
	return nil
}

// validateCounterValue проверяет корректность значения для типа Counter.
func ValidateCounterValue(value string) *types.APIErrorResponse {
	if _, err := strconv.ParseInt(value, 10, 64); err != nil {
		return &types.APIErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid counter value",
		}
	}
	return nil
}

// validateGaugeValue проверяет корректность значения для типа Gauge.
func ValidateGaugeValue(value string) *types.APIErrorResponse {
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return &types.APIErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid gauge value",
		}
	}
	return nil
}

// ValidateMetricValue проверяет значение метрики на корректность.
func ValidateMetricValue(value string, metricType types.MetricType) *types.APIErrorResponse {
	switch metricType {
	case types.Counter:
		return ValidateCounterValue(value)
	case types.Gauge:
		return ValidateGaugeValue(value)
	default:
		return &types.APIErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "unsupported metric type",
		}
	}
}
