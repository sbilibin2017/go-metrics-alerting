package numberprocessor

// Интерфейс с методами Parse и Format.
type NumberProcessorInterface[T int64 | float64] interface {
	Parse(value string) (T, error) // Парсит строку в число
	Format(value T) string         // Форматирует число в строку
}
