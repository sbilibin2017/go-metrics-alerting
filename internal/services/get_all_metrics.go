package services

import (
	"go-metrics-alerting/internal/domain"
)

// Ranger интерфейс для перебора элементов в хранилище.
type Ranger[K comparable, V any] interface {
	Range(callback func(key K, value V) bool)
}

// GetAllMetricsServiceImpl структура для реализации сервиса получения всех метрик.
type GetAllMetricsService struct {
	ranger Ranger[string, *domain.Metrics] // используется для перебора метрик
}

// NewGetAllMetricsService конструктор для GetAllMetricsServiceImpl.
func NewGetAllMetricsService(ranger Ranger[string, *domain.Metrics]) *GetAllMetricsService {
	return &GetAllMetricsService{
		ranger: ranger,
	}
}

// GetAllMetrics получает все метрики из хранилища с помощью Ranger.
func (s *GetAllMetricsService) GetAllMetrics() []*domain.Metrics {
	var metrics []*domain.Metrics
	s.ranger.Range(func(key string, metric *domain.Metrics) bool {
		metrics = append(metrics, metric)
		return true
	})
	return metrics
}
