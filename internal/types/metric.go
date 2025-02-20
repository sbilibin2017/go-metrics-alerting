package types

import (
	"errors"
	"fmt"
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/validators"
)

// UpdateMetricBodyRequest структура для обновления метрики.
type UpdateMetricBodyRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

var ErrDeltaValueConflict error = errors.New("delta and value cannot be set at the same time")

// Validate проверяет поля структуры UpdateMetricBodyRequest.
func (r *UpdateMetricBodyRequest) Validate() error {
	if err := validators.ValidateEmptyString(r.ID); err != nil {
		return err
	}
	if err := validators.ValidateMType(r.MType); err != nil {
		return err
	}
	if r.Delta != nil && r.Value != nil {
		return ErrDeltaValueConflict
	}
	if err := validators.ValidateDelta(r.MType, r.Delta); err != nil {
		return err
	}
	if err := validators.ValidateValue(r.MType, r.Value); err != nil {
		return err
	}
	return nil
}

// ToMetric преобразует запрос в объект domain.Metric.
func (r *UpdateMetricBodyRequest) ToMetric() *domain.Metric {
	switch domain.MType(r.MType) {
	case domain.Counter:
		if r.Delta != nil {
			return &domain.Metric{
				ID:    r.ID,
				MType: domain.Counter,
				Value: fmt.Sprintf("%d", *r.Delta),
			}
		}
	case domain.Gauge:
		if r.Value != nil {
			return &domain.Metric{
				ID:    r.ID,
				MType: domain.Gauge,
				Value: fmt.Sprintf("%f", *r.Value),
			}
		}
	}
	return nil
}

// UpdateMetricPathRequest структура для запроса метрики.
type UpdateMetricPathRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Validate проверяет поля структуры UpdateMetricPathRequest.
func (r *UpdateMetricPathRequest) Validate() error {
	if err := validators.ValidateEmptyString(r.Name); err != nil {
		return err
	}
	if err := validators.ValidateEmptyString(r.Type); err != nil {
		return err
	}
	if err := validators.ValidateEmptyString(r.Value); err != nil {
		return err
	}
	return nil
}

func (r *UpdateMetricPathRequest) ToMetric() *domain.Metric {
	mType := domain.MType(r.Type)
	if mType != domain.Counter && mType != domain.Gauge {
		return nil
	}
	return &domain.Metric{
		ID:    r.Name,
		MType: mType,
		Value: r.Value,
	}
}

// GetMetricBodyRequest структура для запроса метрик.
type GetMetricBodyRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

// Validate проверяет поля структуры GetMetricBodyRequest.
func (r *GetMetricBodyRequest) Validate() error {
	if err := validators.ValidateEmptyString(r.ID); err != nil {
		return err
	}
	if err := validators.ValidateMType(r.MType); err != nil {
		return err
	}
	return nil
}

// ToMetric преобразует запрос в объект domain.Metric.
func (r *GetMetricBodyRequest) ToMetric(value string) *domain.Metric {
	return &domain.Metric{
		ID:    r.ID,
		MType: domain.MType(r.MType),
		Value: value,
	}
}

// GetMetricBodyResponse структура для ответа с метриками.
type GetMetricBodyResponse struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Value string `json:"value"`
}

// ToMetric преобразует ответ в объект domain.Metric.
func (r *GetMetricBodyResponse) ToMetric() *domain.Metric {
	return &domain.Metric{
		ID:    r.ID,
		MType: domain.MType(r.MType),
		Value: r.Value,
	}
}
