package services

import (
	"context"
	"go-metrics-alerting/internal/types"
)

// Интерфейс для получения всех значений метрик
type GetAllRepo interface {
	GetAll(ctx context.Context) [][]string
}

// Сервис для получения всех метрик
type GetAllMetricValuesService struct {
	MetricRepository GetAllRepo
}

// GetAllMetricValues возвращает список всех метрик
func (s *GetAllMetricValuesService) GetAllMetricValues(ctx context.Context) []*types.MetricResponse {
	// Получаем все метрики через репозиторий
	metricsList := s.MetricRepository.GetAll(ctx) // передаем context в репозиторий
	var metrics []*types.MetricResponse

	// Если metricsList пуст, возвращаем пустой срез
	if len(metricsList) == 0 {
		return []*types.MetricResponse{} // Возвращаем пустой срез, а не nil
	}

	// Иначе продолжаем обработку
	for _, metric := range metricsList {
		// Преобразуем строку типа в MetricType
		metricType := types.MetricType(metric[0])

		// Создаем MetricResponse с полями, используя встраивание
		metrics = append(
			metrics,
			&types.MetricResponse{
				UpdateMetricValueRequest: types.UpdateMetricValueRequest{
					Type:  string(metricType),
					Name:  metric[1],
					Value: metric[2],
				},
			},
		)
	}

	return metrics
}
