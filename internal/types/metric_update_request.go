package types

// MetricUpdateRequest определяет структуру запроса для обновления значения метрики.
type MetricUpdateRequest struct {
	Type  MetricType `json:"type"`
	Name  string     `json:"name"`
	Value string     `json:"value"`
}
