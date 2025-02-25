package types

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/validators"
	"net/http"
	"strconv"
)

// Структура для обновления метрики через путь запроса.
type UpdateMetricPathRequest struct {
	ID    string
	MType string
	Value string
}

// Метод преобразования в доменную модель.
func (r *UpdateMetricPathRequest) ToDomain() *domain.Metrics {
	mtype := domain.MetricType(r.MType)
	var value *float64
	var delta *int64
	switch mtype {
	case domain.Counter:
		v, _ := strconv.ParseInt(r.Value, 10, 64)
		delta = &v
	case domain.Gauge:
		v, _ := strconv.ParseFloat(r.Value, 64)
		value = &v
	}
	return &domain.Metrics{
		ID:    r.ID,
		MType: mtype,
		Delta: delta,
		Value: value,
	}
}

// Метод валидации для UpdateMetricPathRequest.
func (r *UpdateMetricPathRequest) Validate() *APIError {
	// Валидация ID
	if err := validators.ValidateEmptyString(r.ID, "ID"); err != nil {
		return &APIError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	// Валидация MType
	if err := validators.ValidateEmptyString(r.MType, "MType"); err != nil {
		return &APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	if err := validators.ValidateMetricType(r.MType); err != nil {
		return &APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// Валидация Value
	if err := validators.ValidateEmptyString(r.Value, "Value"); err != nil {
		return &APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}
