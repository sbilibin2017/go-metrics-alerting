package types

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

// Структура для запроса получения значения метрики
type GetMetricValueRequest struct {
	Type string
	Name string
}

// Структура запроса обновления метрики
type UpdateMetricValueRequest struct {
	Type  string
	Name  string
	Value string
}

type MetricResponse struct {
	UpdateMetricValueRequest
}
