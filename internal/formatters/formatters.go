package formatters

import (
	"errors"
	"fmt"
	"strconv"
)

// Ошибки для парсинга
var (
	ErrParseInt64   = errors.New("failed to parse int64")
	ErrParseFloat64 = errors.New("failed to parse float64")
)

// Реализация парсера и форматтера для int64
type Int64Formatter struct{}

func (h *Int64Formatter) Parse(value string) (int64, error) {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrParseInt64, err)
	}
	return parsedValue, nil
}

func (h *Int64Formatter) Format(value int64) string {
	return strconv.FormatInt(value, 10)
}

// Реализация парсера и форматтера для float64
type Float64Formatter struct{}

func (h *Float64Formatter) Parse(value string) (float64, error) {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrParseFloat64, err)
	}
	return parsedValue, nil
}

func (h *Float64Formatter) Format(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
