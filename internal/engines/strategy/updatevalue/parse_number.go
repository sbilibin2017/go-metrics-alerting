package updatevalue

import "strconv"

// Дженерик парсер для int64 и float64 значений.
func parseNumber[T int64 | float64](value string) (T, error) {
	var zero T
	switch any(zero).(type) {
	case int64:
		parsed, err := strconv.ParseInt(value, 10, 64)
		return T(parsed), err
	case float64:
		parsed, err := strconv.ParseFloat(value, 64)
		return T(parsed), err
	default:
		return zero, strconv.ErrSyntax
	}
}

// FormatNumber форматирует число в строку.
func formatNumber[T int64 | float64](value T) string {
	switch any(value).(type) {
	case int64:
		return strconv.FormatInt(any(value).(int64), 10)
	case float64:
		return strconv.FormatFloat(any(value).(float64), 'f', -1, 64)
	default:
		return ""
	}
}
