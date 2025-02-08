package services

import "go-metrics-alerting/internal/metric/handlers"

// Интерфейс для получения всех значений метрик
type GetAllRepo interface {
	GetAll() [][3]string
}

// Сервис для получения всех метрик
type GetAllMetricValuesService struct {
	metricRepository GetAllRepo
}

// Новый сервис для получения всех метрик
func NewGetAllMetricsService(metricRepository GetAllRepo) *GetAllMetricValuesService {
	return &GetAllMetricValuesService{
		metricRepository: metricRepository,
	}
}

// GetAllMetricValues возвращает список всех метрик
func (s *GetAllMetricValuesService) GetAllMetricValues() []*handlers.MetricResponse {
	metricsList := s.metricRepository.GetAll()
	var metrics []*handlers.MetricResponse

	// Если metricsList пуст, возвращаем пустой срез
	if len(metricsList) == 0 {
		return metrics
	}

	// Иначе продолжаем обработку
	for _, metric := range metricsList {
		// Добавляем метрику в результат
		metrics = append(
			metrics,
			&handlers.MetricResponse{
				Type:  metric[0],
				Name:  metric[1],
				Value: metric[2],
			},
		)
	}

	return metrics
}
