package formatters

import "strconv"

// Форматирует значение типа float64 в строку
func FormatFloat64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// Парсит строковое значение в float64
func ParseFloat64(value string) (float64, bool) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, false
	}
	return v, true
}
