package validators

import (
	"errors"
	"go-metrics-alerting/internal/domain"
	"strings"
)

const (
	EmptyString string = ""
)

// Ошибки, связанные с валидацией
var (
	ErrValueCannotBeEmpty      = errors.New("value cannot be empty")
	ErrInvalidMetricType       = errors.New("invalid metric type")
	ErrDeltaRequiredForCounter = errors.New("delta is required for counter metric type")
	ErrValueRequiredForGauge   = errors.New("value is required for gauge metric type")
)

// ValidateEmptyString проверяет, что строка не пустая.
func ValidateEmptyString(value string) error {
	if strings.TrimSpace(value) == EmptyString {
		return ErrValueCannotBeEmpty
	}
	return nil
}

// ValidateMType проверяет, что mtype имеет допустимое значение.
func ValidateMType(mtype string) error {
	validTypes := []string{string(domain.Counter), string(domain.Gauge)}
	for _, validType := range validTypes {
		if mtype == validType {
			return nil
		}
	}
	return ErrInvalidMetricType
}

// ValidateDelta проверяет, что Delta указано, если тип метрики "counter".
func ValidateDelta(mtype string, delta *int64) error {
	if mtype == string(domain.Counter) && delta == nil {
		return ErrDeltaRequiredForCounter
	}
	return nil
}

// ValidateValue проверяет, что Value указано, если тип метрики "gauge".
func ValidateValue(mtype string, value *float64) error {
	if mtype == string(domain.Gauge) && value == nil {
		return ErrValueRequiredForGauge
	}
	return nil
}
