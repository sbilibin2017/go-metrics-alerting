package domain

// Тип для представления типов метрик
type MetricType string

// Возможные значения для типа метрики
const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

// Структура для метрик
type Metrics struct {
	ID    string
	MType MetricType
	Delta *int64
	Value *float64
}
