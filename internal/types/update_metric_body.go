package types

import (
	"fmt"
	"go-metrics-alerting/internal/apierror"
	"go-metrics-alerting/internal/domain"
	"net/http"
	"regexp"
)

//go:generate easyjson -all update_metric_body.go
type UpdateMetricBodyRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// ValidateUpdateMetricBodyRequest проверяет данные, полученные из тела запроса
func (r *UpdateMetricBodyRequest) Validate() *apierror.APIError {
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

	// Для gauge значение должно быть числом с плавающей точкой
	if r.Value != nil && r.Delta != nil {
		return &apierror.APIError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("only value('%v') or delta('%v')) must be set", r.Value, r.Delta),
		}
	}

	// Для gauge значение должно быть числом с плавающей точкой
	if r.MType == string(domain.Gauge) && r.Value == nil {
		return &apierror.APIError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("value is required for gauge metrics, received: '%v'", r.Value),
		}
	}

	// Для counter значение должно быть целым числом
	if r.MType == string(domain.Counter) && r.Delta == nil {
		return &apierror.APIError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("delta is required for counter metrics, received: '%v'", r.Delta),
		}
	}

	return nil
}

func (r *UpdateMetricBodyRequest) ToDomain() *domain.Metrics {
	return &domain.Metrics{
		ID:    r.ID,
		MType: domain.MType(r.MType),
		Delta: r.Delta,
		Value: r.Value,
	}
}
