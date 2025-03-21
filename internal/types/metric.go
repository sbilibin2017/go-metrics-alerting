package types

type MType string

const (
	Gauge   MType = "gauge"
	Counter MType = "counter"
)

type MetricID struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Metrics struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
