package numberprocessor

import "strconv"

// Реализация для int64.
type Int64Processor struct{}

// Парсинг для int64.
func (Int64Processor) Parse(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// Форматирование для int64.
func (Int64Processor) Format(value int64) string {
	return strconv.FormatInt(value, 10)
}
