package types

// Определяем enum для типов метрик
type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)
