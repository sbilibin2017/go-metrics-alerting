package numberprocessor

import "strconv"

// Реализация для int64.
type Int64ProcessorEngine struct{}

func NewInt64ProcessorEngine() *Int64ProcessorEngine {
	return &Int64ProcessorEngine{}
}

// Парсинг для int64.
func (Int64ProcessorEngine) Parse(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// Форматирование для int64.
func (Int64ProcessorEngine) Format(value int64) string {
	return strconv.FormatInt(value, 10)
}
