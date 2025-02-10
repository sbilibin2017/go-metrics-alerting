package types

// UpdateMetricValueRequest представляет запрос на обновление метрики.
type UpdateMetricValueRequest struct {
	Type  MetricType
	Name  string
	Value string
}
