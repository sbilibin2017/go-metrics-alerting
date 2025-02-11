package types

// Определяем enum для типов метрик
type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

// UpdateMetricValueRequest представляет запрос на обновление метрики.
type UpdateMetricValueRequest struct {
	Type  MetricType
	Name  string
	Value string
}

// GetMetricValueRequest представляет запрос на получение значения метрики.
type GetMetricValueRequest struct {
	Type MetricType
	Name string
}

// MetricResponse представляет ответ.
type MetricResponse struct {
	Name  string
	Value string
}
