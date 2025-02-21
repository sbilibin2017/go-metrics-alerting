package types

// MetricType является типом для определения типа метрики (gauge или counter).
type MType string

// Возможные значения для MetricType
const (
	Gauge   MType = "gauge"
	Counter MType = "counter"
)

// Metrics является структурой для метрик
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType MType    `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
