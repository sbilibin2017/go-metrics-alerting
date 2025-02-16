package validators

import (
	"go-metrics-alerting/internal/types"
	"net/http"
)

// Константы сообщений об ошибках
const (
	ErrIDRequired        = "ID is required"
	ErrMTypeRequired     = "Metric type is required"
	ErrInvalidMetricType = "Invalid metric type, must be 'counter' or 'gauge'"
	ErrDeltaRequired     = "Delta is required for 'counter' type"
	ErrValueRequired     = "Value is required for 'gauge' type"
)

// ValidateMetrics выполняет валидацию структуры Metrics
func ValidateMetrics(m *types.MetricsRequest) *types.APIErrorResponse {
	if isIDEmpty(m.ID) {
		return types.NewAPIErrorResponse(http.StatusNotFound, ErrIDRequired)
	}
	if isMTypeEmpty(m.MType) {
		return types.NewAPIErrorResponse(http.StatusNotFound, ErrMTypeRequired)
	}
	if isMTypeInvalid(m.MType) {
		return types.NewAPIErrorResponse(http.StatusBadRequest, ErrInvalidMetricType)
	}
	if isDeltaEmpty(m.MType, m.Delta) {
		return types.NewAPIErrorResponse(http.StatusBadRequest, ErrDeltaRequired)
	}
	if isValueEmpty(m.MType, m.Value) {
		return types.NewAPIErrorResponse(http.StatusBadRequest, ErrValueRequired)
	}
	return nil
}

// ValidateMetricValueRequest выполняет валидацию структуры MetricValueRequest
func ValidateMetricValueRequest(m *types.MetricValueRequest) *types.APIErrorResponse {
	if isIDEmpty(m.ID) {
		return types.NewAPIErrorResponse(http.StatusNotFound, ErrIDRequired)
	}
	if isMTypeEmpty(m.MType) {
		return types.NewAPIErrorResponse(http.StatusNotFound, ErrMTypeRequired)
	}
	if isMTypeInvalid(m.MType) {
		return types.NewAPIErrorResponse(http.StatusBadRequest, ErrInvalidMetricType)
	}
	return nil
}

// Проверяет, что ID не пустой
func isIDEmpty(id string) bool {
	return id == types.EmptyString
}

// Проверяет, что MType не пустой
func isMTypeEmpty(mType types.MType) bool {
	return string(mType) == types.EmptyString
}

// Проверяет, что MType имеет корректное значение
func isMTypeInvalid(mType types.MType) bool {
	return mType != types.Counter && mType != types.Gauge
}

// Проверяет, что Delta указана для типа Counter
func isDeltaEmpty(mtype types.MType, delta *int64) bool {
	return mtype == types.Counter && delta == nil
}

// Проверяет, что Value указан для типа Gauge
func isValueEmpty(mtype types.MType, value *float64) bool {
	return mtype == types.Gauge && value == nil
}
