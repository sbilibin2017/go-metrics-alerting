package types

// MType определяет тип метрики.
type MType string

const (
	Counter MType = "counter"
	Gauge   MType = "gauge"
)

// MetricsRequest используется для обновления метрики.
type MetricsRequest struct {
	ID    string   `json:"id"`
	MType MType    `json:"mtype"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricsResponse struct {
	MetricsRequest
}

// MetricValueRequest используется для запроса значения метрики.
type MetricValueRequest struct {
	ID    string `json:"id"`
	MType MType  `json:"mtype"`
}

// MetricValueRequest используется для запроса значения метрики.
type MetricValueResponse struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
