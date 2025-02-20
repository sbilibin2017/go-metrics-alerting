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

// Форматирует значение типа int64 в строку
func FormatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}

// Парсит строковое значение в int64
func ParseInt64(value string) (int64, bool) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
