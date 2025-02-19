package validators

import (
	"errors"
	"go-metrics-alerting/internal/types"
	"strconv"
)

// Ошибки валидации
var (
	ErrEmptyID           = errors.New("id cannot be empty")
	ErrInvalidMType      = errors.New("invalid metric type")
	ErrInvalidDelta      = errors.New("counter metric must have a delta value")
	ErrInvalidValue      = errors.New("gauge metric must have a value")
	ErrInvalidCounterVal = errors.New("counter value must be a valid integer")
	ErrInvalidGaugeVal   = errors.New("gauge value must be a valid float")
)

// ValidateEmptyString валидатор для проверки пустой строки
type ValidateEmptyString struct{}

func (v *ValidateEmptyString) Validate(id string) error {
	if id == "" {
		return ErrEmptyID
	}
	return nil
}

// ValidateMType валидатор для проверки типа метрики
type ValidateMType struct{}

func (v *ValidateMType) Validate(mType types.MType) error {
	if mType != types.Counter && mType != types.Gauge {
		return ErrInvalidMType
	}
	return nil
}

// ValidateDelta валидатор для проверки Delta для счетчиков
type ValidateDelta struct{}

func (v *ValidateDelta) Validate(mType types.MType, delta *int64) error {
	if mType == types.Counter && delta == nil {
		return ErrInvalidDelta
	}
	return nil
}

// ValidateValue валидатор для проверки Value для Gauge
type ValidateValue struct{}

func (v *ValidateValue) Validate(mType types.MType, value *float64) error {
	if mType == types.Gauge && value == nil {
		return ErrInvalidValue
	}
	return nil
}

// ValidateCounterValue валидатор для проверки значения для счетчика
type ValidateCounterValue struct{}

func (v *ValidateCounterValue) Validate(value string) error {
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return ErrInvalidCounterVal
	}
	return nil
}

// ValidateGaugeValue валидатор для проверки значения для Gauge
type ValidateGaugeValue struct{}

func (v *ValidateGaugeValue) Validate(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return ErrInvalidGaugeVal
	}
	return nil
}
