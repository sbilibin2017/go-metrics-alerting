package formatters

import (
	"errors"
	"fmt"
	"strconv"
)

// Определяем ошибку для некорректного парсинга float64
var ErrParseFloat64 = errors.New("failed to parse float64")

// ParseFloat64 парсит строковое представление числа в float64.
func ParseFloat64(value string) (float64, error) {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrParseFloat64, err)
	}
	return parsedValue, nil
}

// FormatFloat64 форматирует число типа float64 в строку.
func FormatFloat64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
