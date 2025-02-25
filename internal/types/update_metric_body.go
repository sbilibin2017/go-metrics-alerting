package types

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/validators"
	"net/http"
)

//go:generate easyjson -all update_metric_body.go
type UpdateMetricBodyRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

//go:generate easyjson -all update_metric_body.go
type UpdateMetricBodyResponse struct {
	UpdateMetricBodyRequest
}

// Метод преобразования в доменную модель.
func (r *UpdateMetricBodyRequest) ToDomain() *domain.Metrics {
	return &domain.Metrics{
		ID:    r.ID,
		MType: domain.MetricType(r.MType),
		Delta: r.Delta,
		Value: r.Value,
	}
}

// Метод валидации для UpdateMetricBodyRequest.
func (r *UpdateMetricBodyRequest) Validate() *APIError {
	// Валидация ID
	if err := validators.ValidateEmptyString(r.ID, "ID"); err != nil {
		return &APIError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	// Валидация MType
	if err := validators.ValidateMetricType(r.MType); err != nil {
		return &APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// Валидация Delta/Value в зависимости от типа метрики
	switch r.MType {
	case string(domain.Counter):
		if err := validators.ValidateInt64Ptr(r.Delta, "Delta"); err != nil {
			return &APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}
	case string(domain.Gauge):
		if err := validators.ValidateFloat64Ptr(r.Value, "Value"); err != nil {
			return &APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}
	}

	return nil
}
