package types

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

func GetAllMetricTypes() []MetricType {
	return []MetricType{Counter, Gauge}
}

type MetricID struct {
	ID   string     `json:"id"`
	Type MetricType `json:"type"`
}

type Metrics struct {
	MetricID
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
