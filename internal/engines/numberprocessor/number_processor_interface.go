package numberprocessor

// Интерфейс для парсинга и форматирования значений.
type NumberProcessorEngineInterface[T int64 | float64] interface {
	// Parse парсит строковое значение в числовое.
	Parse(value string) (T, error)

	// Format форматирует числовое значение в строку.
	Format(value T) string
}
