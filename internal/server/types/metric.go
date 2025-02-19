package types

type MType string

const (
	Counter MType = "counter"
	Gauge   MType = "gauge"
)

type UpdateMetricsRequest struct {
	ID    string   `json:"id"`
	MType MType    `json:"mtype"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type UpdateMetricsResponse struct {
	UpdateMetricsRequest
}

type GetMetricValueRequest struct {
	ID    string `json:"id"`
	MType MType  `json:"mtype"`
}

type GetMetricValueResponse struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
