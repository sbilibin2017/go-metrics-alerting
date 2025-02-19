package formatters

import (
	"errors"
	"fmt"
	"strconv"
)

// Определяем ошибку для некорректного парсинга int64
var ErrParseInt64 = errors.New("failed to parse int64")

// ParseInt64 парсит строковое представление числа в int64.
func ParseInt64(value string) (int64, error) {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrParseInt64, err)
	}
	return parsedValue, nil
}

// FormatInt64 форматирует число типа int64 в строку.
func FormatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}
