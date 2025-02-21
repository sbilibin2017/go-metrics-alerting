package types

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/validators" // Импортируем пакет с функциями валидации
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
func (r *UpdateMetricBodyRequest) Validate() error {
	if err := validators.ValidateString(r.ID); err != nil {
		return err
	}
	mType := domain.MType(r.MType)
	if err := validators.ValidateMType(mType); err != nil {
		return err
	}
	switch mType {
	case domain.Counter:
		if err := validators.ValidateDelta(mType, r.Delta); err != nil {
			return err
		}
	case domain.Gauge:
		if err := validators.ValidateValue(mType, r.Value); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMetricPathRequest является структурой для обновления метрики через путь запроса.
type UpdateMetricPathRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
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
func (r *UpdateMetricPathRequest) Validate() error {
	if err := validators.ValidateString(r.ID); err != nil {
		return err
	}
	if err := validators.ValidateMType(domain.MType(r.MType)); err != nil {
		return err
	}
	if err := validators.ValidateValueString(domain.MType(r.MType), r.Value); err != nil {
		return err
	}
	return nil
}

// GetMetricBodyRequest является структурой для получения метрики через тело запроса.
type GetMetricBodyRequest struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

// Метод для валидации GetMetricBodyRequest
func (r *GetMetricBodyRequest) Validate() error {
	if err := validators.ValidateString(r.ID); err != nil {
		return err
	}
	if err := validators.ValidateMType(domain.MType(r.MType)); err != nil {
		return err
	}
	return nil
}

// GetMetricPathRequest является структурой для получения метрики через путь запроса.
type GetMetricPathRequest struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

// Метод для валидации GetMetricPathRequest
func (r *GetMetricPathRequest) Validate() error {
	if err := validators.ValidateString(r.ID); err != nil {
		return err
	}
	if err := validators.ValidateMType(domain.MType(r.MType)); err != nil {
		return err
	}
	return nil
}
