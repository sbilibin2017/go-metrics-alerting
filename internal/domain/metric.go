package domain

// Создаем тип MType как строку
type MType string

// Определяем константы для возможных значений MType
const (
	Gauge   MType = "gauge"   // runtime метрика
	Counter MType = "counter" // счетчик
)

// Структура Metric с типом MType
type Metric struct {
	ID    string // идентификатор метрики
	MType MType  // тип метрики
	Value string // значение метрики
}
