package types

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/validators"
	"net/http"
	"strconv"
)

//go:generate easyjson -all get_metric.go
type GetMetricRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

// Метод валидации для GetMetricRequest.
func (r *GetMetricRequest) Validate() *APIError {
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

	// Валидация правильности типа MType
	if err := validators.ValidateMetricType(r.MType); err != nil {
		return &APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}

// Структура для получения метрики через запрос.
type GetMetricPathResponse struct {
	Value string `json:"value"`
}

// Метод преобразования из доменной модели в структуру GetMetricPathResponse.
func (r *GetMetricPathResponse) FromDomain(metric *domain.Metrics) GetMetricPathResponse {
	var valueStr string
	switch metric.MType {
	case domain.Counter:
		valueStr = strconv.FormatInt(*metric.Delta, 10)
	case domain.Gauge:
		valueStr = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	}
	return GetMetricPathResponse{
		Value: valueStr,
	}
}

//go:generate easyjson -all get_metric.go
type GetMetricBodyResponse struct {
	GetMetricRequest
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
