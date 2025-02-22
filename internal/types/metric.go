package types

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"strconv"
)

// UpdateMetricBodyRequest является структурой для обновления метрики через тело запроса.
type UpdateMetricBodyRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// Метод для конвертации UpdateMetricBodyRequest в доменную сущность Metrics
func (r *UpdateMetricBodyRequest) ToDomain() *domain.Metrics {
	return &domain.Metrics{
		ID:    r.ID,
		MType: domain.MType(r.MType),
		Delta: r.Delta,
		Value: r.Value,
	}
}

// Метод для валидации UpdateMetricBodyRequest
func (r *UpdateMetricBodyRequest) Validate() *APIError {
	// Валидация ID
	if r.ID == "" {
		return &APIError{Status: http.StatusNotFound, Message: "id is required"}
	}

	// Валидация типа метрики
	mType := domain.MType(r.MType)
	if mType != domain.Counter && mType != domain.Gauge {
		return &APIError{Status: http.StatusBadRequest, Message: "invalid metric type"}
	}

	// Проверка на необходимое значение для Delta или Value
	if (r.Delta != nil && r.Value != nil) || (r.Delta == nil && r.Value == nil) {
		return &APIError{Status: http.StatusBadRequest, Message: "delta or value must be nil"}
	}

	// Валидация для типа Counter (требуется delta)
	if mType == domain.Counter && r.Delta == nil {
		return &APIError{Status: http.StatusBadRequest, Message: "delta is required for Counter metric"}
	}

	// Валидация для типа Gauge (требуется value)
	if mType == domain.Gauge && r.Value == nil {
		return &APIError{Status: http.StatusBadRequest, Message: "value is required for Gauge metric"}
	}

	return nil
}

// UpdateMetricPathRequest является структурой для обновления метрики через путь запроса.
type UpdateMetricPathRequest struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
	Value string `json:"value"`
}

// Метод для конвертации UpdateMetricPathRequest в доменную сущность Metrics
func (r *UpdateMetricPathRequest) ToDomain() *domain.Metrics {
	var delta *int64
	var value *float64
	mType := domain.MType(r.MType)
	switch mType {
	case domain.Gauge:
		parsedValue, err := strconv.ParseFloat(r.Value, 64)
		if err == nil {
			value = &parsedValue
		}
	case domain.Counter:
		parsedDelta, err := strconv.ParseInt(r.Value, 10, 64)
		if err == nil {
			delta = &parsedDelta
		}
	}
	return &domain.Metrics{
		ID:    r.ID,
		MType: mType,
		Delta: delta,
		Value: value,
	}
}

// Метод для валидации UpdateMetricPathRequest
func (r *UpdateMetricPathRequest) Validate() *APIError {
	// Валидация ID
	if r.ID == "" {
		return &APIError{Status: http.StatusNotFound, Message: "id is required"}
	}

	// Валидация типа метрики
	if r.MType != string(domain.Counter) && r.MType != string(domain.Gauge) {
		return &APIError{Status: http.StatusBadRequest, Message: "invalid metric type"}
	}

	// Валидация значения строки
	mType := domain.MType(r.MType)
	switch mType {
	case domain.Gauge:
		_, err := strconv.ParseFloat(r.Value, 64)
		if err != nil {
			return &APIError{Status: http.StatusBadRequest, Message: "invalid value for Gauge metric, must be a valid float"}
		}
	case domain.Counter:
		_, err := strconv.ParseInt(r.Value, 10, 64)
		if err != nil {
			return &APIError{Status: http.StatusBadRequest, Message: "invalid value for Counter metric, must be a valid integer"}
		}
	}

	return nil
}

// GetMetricBodyRequest является структурой для получения метрики через тело запроса.
type GetMetricRequest struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

// Метод для валидации GetMetricBodyRequest
func (r *GetMetricRequest) Validate() *APIError {
	// Валидация ID
	if r.ID == "" {
		return &APIError{Status: http.StatusNotFound, Message: "id is required"}
	}

	// Валидация типа метрики
	if r.MType != "counter" && r.MType != "gauge" {
		return &APIError{Status: http.StatusBadRequest, Message: "invalid metric type"}
	}

	return nil
}
