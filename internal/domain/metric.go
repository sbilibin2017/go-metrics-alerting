package domain

// MetricType является типом для определения типа метрики (gauge или counter).
type MType string

// Возможные значения для MetricType
const (
	Gauge   MType = "gauge"
	Counter MType = "counter"
)

// Metrics является структурой для метрик
type Metrics struct {
	ID    string
	MType MType
	Delta *int64
	Value *float64
}
