package validators

import (
	"errors"
	"go-metrics-alerting/internal/api/types"
)

const (
	EmptyString string = ""
)

// Определяем ошибки валидации как константы.
var (
	ErrEmptyID      = errors.New("id cannot be empty")
	ErrInvalidMType = errors.New("invalid metric type")
	ErrInvalidDelta = errors.New("counter metric must have a delta value")
	ErrInvalidValue = errors.New("gauge metric must have a value")
)

// ValidateID проверяет, что ID не пустой.
func ValidateID(id string) error {
	if id == EmptyString {
		return ErrEmptyID
	}
	return nil
}

// ValidateMType проверяет, является ли тип метрики допустимым.
func ValidateMType(mType types.MType) error {
	if mType != types.Counter && mType != types.Gauge {
		return ErrInvalidMType
	}
	return nil
}

// ValidateDelta проверяет, что Delta задана для счетчиков.
func ValidateDelta(mType types.MType, delta *int64) error {
	if mType == types.Counter && delta == nil {
		return ErrInvalidDelta
	}
	return nil
}

// ValidateValue проверяет, что Value задано для Gauge.
func ValidateValue(mType types.MType, value *float64) error {
	if mType == types.Gauge && value == nil {
		return ErrInvalidValue
	}
	return nil
}
