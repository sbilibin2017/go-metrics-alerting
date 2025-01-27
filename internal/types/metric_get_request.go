package types

// MetricGetRequest определяет структуру запроса для получения значения метрики.
type MetricGetRequest struct {
	Type MetricType `json:"type"`
	Name string     `json:"name"`
}
