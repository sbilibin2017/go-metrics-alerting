package validators

import (
	"errors"
	"go-metrics-alerting/internal/domain"
	"strconv"
)

// Функция для валидации строки
func ValidateString(s string) error {
	if s == "" {
		return errors.New("id is required")
	}
	return nil
}

// Функция для валидации типа метрики
func ValidateMType(mType domain.MType) error {
	if mType != domain.Counter && mType != domain.Gauge {
		return errors.New("invalid metric type")
	}
	return nil
}

// Функция для валидации Delta (для Counter)
func ValidateDelta(mType domain.MType, delta *int64) error {
	if mType == domain.Counter && delta == nil {
		return errors.New("delta is required for Counter metric")
	}
	return nil
}

// Функция для валидации Value (для Gauge)
func ValidateValue(mType domain.MType, value *float64) error {
	if mType == domain.Gauge && value == nil {
		return errors.New("value is required for Gauge metric")
	}
	return nil
}

// Функция для валидации значения метрики в формате строки
func ValidateValueString(mType domain.MType, value string) error {
	if err := ValidateString(value); err != nil {
		return err
	}
	switch mType {
	case domain.Gauge:
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("invalid value for Gauge metric, must be a valid float")
		}
	case domain.Counter:
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errors.New("invalid value for Counter metric, must be a valid integer")
		}
	default:
		return errors.New("unsupported metric type")
	}
	return nil
}
