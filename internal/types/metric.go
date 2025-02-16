package types

type MType string

const (
	Counter MType = "counter"
	Gauge   MType = "gauge"
)

type MetricsRequest struct {
	ID    string   `json:"id"`
	MType MType    `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricValueRequest struct {
	ID    string `json:"id"`
	MType MType  `json:"type"`
}
