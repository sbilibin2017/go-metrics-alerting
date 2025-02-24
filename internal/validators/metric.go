package validators

import (
	"fmt"
	"go-metrics-alerting/internal/domain"
	"strings"
)

// ValidateEmptyString проверяет, что строка не пуста.
func ValidateEmptyString(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateMetricType проверяет, что тип метрики является допустимым (counter или gauge).
func ValidateMetricType(mtype string) error {
	validTypes := map[string]bool{
		string(domain.Counter): true,
		string(domain.Gauge):   true,
	}

	mtype = strings.ToLower(mtype)
	if !validTypes[mtype] {
		return fmt.Errorf("invalid metric type: %s", mtype)
	}

	return nil
}

// ValidateInt64Ptr проверяет, что указатель на int64 не равен nil.
func ValidateInt64Ptr(value *int64, fieldName string) error {
	if value == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}
	return nil
}

// ValidateFloat64Ptr проверяет, что указатель на float64 не равен nil.
func ValidateFloat64Ptr(value *float64, fieldName string) error {
	if value == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}
	return nil
}
