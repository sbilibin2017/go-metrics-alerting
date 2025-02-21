package types

import (
	"go-metrics-alerting/internal/domain"
	"strconv"
)

// Интерфейсы для валидации

type StringValidator interface {
	Validate(s string) error
}

type MetricTypeValidator interface {
	Validate(mType domain.MType) error
}

type DeltaValidator interface {
	Validate(mType domain.MType, delta *int64) error
}

type ValueValidator interface {
	Validate(mType domain.MType, value *float64) error
}

type ValueStringValidator interface {
	Validate(value string) error
}

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
func (r *UpdateMetricBodyRequest) Validate(
	stringValidator StringValidator,
	metricTypeValidator MetricTypeValidator,
	deltaValidator DeltaValidator,
	valueValidator ValueValidator,
) error {
	if err := stringValidator.Validate(r.ID); err != nil {
		return err
	}
	mType := domain.MType(r.MType)
	if err := metricTypeValidator.Validate(mType); err != nil {
		return err
	}
	switch mType {
	case domain.Counter:
		if err := deltaValidator.Validate(mType, r.Delta); err != nil {
			return err
		}
	case domain.Gauge:
		if err := valueValidator.Validate(mType, r.Value); err != nil {
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
func (r *UpdateMetricPathRequest) Validate(
	stringValidator StringValidator,
	metricTypeValidator MetricTypeValidator,
	valueStringValidator ValueStringValidator,
) error {
	if err := stringValidator.Validate(r.ID); err != nil {
		return err
	}
	mType := domain.MType(r.MType)
	if err := metricTypeValidator.Validate(mType); err != nil {
		return err
	}
	if err := valueStringValidator.Validate(r.Value); err != nil {
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
func (r *GetMetricBodyRequest) Validate(
	stringValidator StringValidator,
	metricTypeValidator MetricTypeValidator,
) error {
	if err := stringValidator.Validate(r.ID); err != nil {
		return err
	}
	mType := domain.MType(r.MType)
	if err := metricTypeValidator.Validate(mType); err != nil {
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
func (r *GetMetricPathRequest) Validate(
	stringValidator StringValidator,
	metricTypeValidator MetricTypeValidator,
) error {
	if err := stringValidator.Validate(r.ID); err != nil {
		return err
	}
	mType := domain.MType(r.MType)
	if err := metricTypeValidator.Validate(mType); err != nil {
		return err
	}
	return nil
}
