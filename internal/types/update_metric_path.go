package types

import (
	"fmt"
	"go-metrics-alerting/internal/apierror"
	"go-metrics-alerting/internal/domain"
	"net/http"
	"regexp"
	"strconv"
)

type UpdateMetricPathRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Value string `json:"value"`
}

// ValidateUpdateMetricPathRequest проверяет параметры пути
func (r UpdateMetricPathRequest) Validate() *apierror.APIError {
	// Проверка на допустимый ID
	if len(r.ID) == 0 {
		return &apierror.APIError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("metric ID must not be empty, received: '%s'", r.ID),
		}
	}

	// Проверка ID с использованием регулярного выражения
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !re.MatchString(r.ID) {
		return &apierror.APIError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("invalid metric ID format, received: '%s'", r.ID),
		}
	}

	// Проверка на допустимый тип метрики
	if r.MType != string(domain.Gauge) && r.MType != string(domain.Counter) {
		return &apierror.APIError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("invalid metric type, received: '%s'", r.MType),
		}
	}

	// Проверка на валидность значения
	if r.MType == string(domain.Gauge) {
		if _, err := strconv.ParseFloat(r.Value, 64); err != nil {
			return &apierror.APIError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid value for gauge, received: '%s'", r.Value),
			}
		}
	} else if r.MType == string(domain.Counter) {
		if _, err := strconv.ParseInt(r.Value, 10, 64); err != nil {
			return &apierror.APIError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid value for counter, received: '%s'", r.Value),
			}
		}
	}

	return nil
}

// ToDomain преобразует UpdateMetricPathRequest в структуру Metrics
func (r *UpdateMetricPathRequest) ToDomain() *domain.Metrics {
	var metric domain.Metrics
	metric.ID = r.ID
	metric.MType = domain.MType(r.MType)
	if r.MType == string(domain.Gauge) {
		value, err := strconv.ParseFloat(r.Value, 64)
		if err == nil {
			metric.Value = &value
		}
	}
	if r.MType == string(domain.Counter) {
		delta, err := strconv.ParseInt(r.Value, 10, 64)
		if err == nil {
			metric.Delta = &delta
		}
	}
	return &metric
}
