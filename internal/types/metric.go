package types

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

func GetAllMetricTypes() []MetricType {
	return []MetricType{
		GaugeType,
		CounterType,
	}
}

type MetricKey struct {
	Type string
	Name string
}

type UpdateMetricRequest struct {
	Type  string
	Name  string
	Value string
}

type GetMetricRequest struct {
	Type string
	Name string
}
