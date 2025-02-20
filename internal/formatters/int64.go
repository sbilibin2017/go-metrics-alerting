package formatters

import "strconv"

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
