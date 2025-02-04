package numberprocessor

import "strconv"

// Float64Processor реализует интерфейс NumberProcessorInterface для типа float64.
type Float64ProcessorEngine struct{}

func NewFloat64ProcessorEngine() *Float64ProcessorEngine {
	return &Float64ProcessorEngine{}
}

// Парсинг для float64.
func (Float64ProcessorEngine) Parse(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// Форматирование для float64.
func (Float64ProcessorEngine) Format(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
