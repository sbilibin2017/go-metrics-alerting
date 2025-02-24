package types

import (
	"fmt"
	"go-metrics-alerting/internal/domain"
)

// Структура для отображения метрик на странице.
type GetAllMetricsResponse struct {
	ID    string
	Value string
}

// Метод для преобразования данных из доменной модели в структуру GetAllMetricsResponse
// и возвращения нового объекта.
func (r *GetAllMetricsResponse) FromDomain(metric *domain.Metrics) GetAllMetricsResponse {
	resp := GetAllMetricsResponse{
		ID: metric.ID,
	}
	switch metric.MType {
	case domain.Counter:
		resp.Value = fmt.Sprintf("%d", *metric.Delta)
	case domain.Gauge:
		resp.Value = fmt.Sprintf("%f", *metric.Value)
	}
	return resp
}
