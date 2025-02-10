package types

// GetMetricValueRequest представляет запрос на получение значения метрики.
type GetMetricValueRequest struct {
	Type string
	Name string
}
