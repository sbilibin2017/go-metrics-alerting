package validators

import (
	"errors"
	"go-metrics-alerting/internal/server/types"
	"strconv"
)

const EmptyString string = ""

// Определяем ошибки валидации как константы.
var (
	ErrEmptyID           = errors.New("id cannot be empty")
	ErrInvalidMType      = errors.New("invalid metric type")
	ErrInvalidDelta      = errors.New("counter metric must have a delta value")
	ErrInvalidValue      = errors.New("gauge metric must have a value")
	ErrInvalidCounterVal = errors.New("counter value must be a valid integer")
	ErrInvalidGaugeVal   = errors.New("gauge value must be a valid float")
)

// ValidateEmptyString проверяет, что ID не пустой.
func ValidateEmptyString(id string) error {
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

// ValidateCounterValue проверяет, что значение для счетчика является допустимым целым числом.
func ValidateCounterValue(value string) error {
	// Преобразуем строку в int64 для счетчика
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return ErrInvalidCounterVal
	}
	return nil
}

// ValidateGaugeValue проверяет, что значение для Gauge является допустимым числом с плавающей точкой.
func ValidateGaugeValue(value string) error {
	// Преобразуем строку в float64 для Gauge
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return ErrInvalidGaugeVal
	}
	return nil
}
