package format

import (
	"strconv"
)

// FormatInt64 форматирует число типа int64 в строку
func FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10) // 10 — это основание системы счисления (десятичная система)
}

// ParseInt64 парсит строку в число типа int64
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64) // 10 — основание системы счисления, 64 — размерность числа
}

// FormatFloat64 форматирует число типа float64 в строку с максимальной точностью
func FormatFloat64(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// ParseFloat64 парсит строку в число типа float64
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
