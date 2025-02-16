package validators

import (
	"errors"
	"go-metrics-alerting/internal/types"
)

const (
	// EmptyString - пустая строка для проверки.
	EmptyString = ""
)

// Переменные для ошибок с использованием errors.New
var (
	ErrIDEmpty      = errors.New("id cannot be empty")
	ErrMTypeEmpty   = errors.New("mtype cannot be empty")
	ErrMTypeInvalid = errors.New("mtype must be either 'counter' or 'gauge'")
	ErrDeltaEmpty   = errors.New("delta must be provided for Counter metrics")
	ErrValueEmpty   = errors.New("value must be provided for Gauge metrics")
)

// validateID проверяет, что ID метрики не пустой.
func ValidateID(id string) error {
	if id == EmptyString {
		return ErrIDEmpty
	}
	return nil
}

// validateMType проверяет, что MType не пустой и имеет корректное значение.
func ValidateMType(mType types.MType) error {
	if string(mType) == EmptyString {
		return ErrMTypeEmpty
	}
	if mType != types.Counter && mType != types.Gauge {
		return ErrMTypeInvalid
	}
	return nil
}

// validateDelta проверяет, что Delta указана для типа Counter.
func ValidateDelta(mtype types.MType, delta *int64) error {
	if mtype == types.Counter && delta == nil {
		return ErrDeltaEmpty
	}
	return nil
}

// validateValue проверяет, что Value указана для типа Gauge.
func ValidateValue(mtype types.MType, value *float64) error {
	if mtype == types.Gauge && value == nil {
		return ErrValueEmpty
	}
	return nil
}
