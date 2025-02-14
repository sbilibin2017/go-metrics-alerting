package types

import (
	"net/http"
	"strconv"
)

const emptyString string = ""

// MetricType - тип метрики
type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

// UpdateMetricValueRequest представляет запрос на обновление метрики.
type UpdateMetricValueRequest struct {
	Type  MetricType
	Name  string
	Value string
}

// validateName проверяет, что имя метрики не пустое.
func validateName(name string) *APIErrorResponse {
	if name == emptyString {
		return &APIErrorResponse{
			Code:    http.StatusNotFound,
			Message: "metric name is required",
		}
	}
	return nil
}

// validateType проверяет, что тип метрики является допустимым.
func validateType(t MetricType) *APIErrorResponse {
	if t != Counter && t != Gauge {
		return &APIErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid metric type",
		}
	}
	return nil
}

// Validate проверяет, что запрос на обновление метрики корректен.
func (req *UpdateMetricValueRequest) Validate() *APIErrorResponse {
	if err := validateName(req.Name); err != nil {
		return err
	}

	if err := validateType(req.Type); err != nil {
		return err
	}

	// Проверка значений для разных типов метрик
	switch req.Type {
	case Counter:
		if _, err := strconv.ParseInt(req.Value, 10, 64); err != nil {
			return &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid counter value",
			}
		}
	case Gauge:
		if _, err := strconv.ParseFloat(req.Value, 64); err != nil {
			return &APIErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid gauge value",
			}
		}
	}

	return nil
}

// GetMetricValueRequest представляет запрос на получение значения метрики.
type GetMetricValueRequest struct {
	Type MetricType
	Name string
}

// Validate проверяет, что запрос на получение метрики корректен.
func (req *GetMetricValueRequest) Validate() *APIErrorResponse {
	if err := validateName(req.Name); err != nil {
		return err
	}

	if err := validateType(req.Type); err != nil {
		return err
	}

	return nil
}

// MetricResponse представляет ответ.
type MetricResponse struct {
	Name  string
	Value string
}
