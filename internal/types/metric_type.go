package types

// MetricType определяет энам для типов метрик.
type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)
