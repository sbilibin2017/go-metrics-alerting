package services

import (
	"errors"
	"go-metrics-alerting/internal/domain"
)

// ErrMetricNotFound is the error returned when a metric is not found in storage.
var ErrMetricNotFound = errors.New("metric not found")

// Getter интерфейс для чтения данных.
type Getter[K comparable, V any] interface {
	Get(key K) (V, bool)
}

// GetMetricPathService интерфейс для получения метрики по ID и типу.
type GetMetricPathService interface {
	GetMetric(id string, mtype domain.MetricType) (*domain.Metrics, error)
}

// GetMetricService структура для реализации сервиса получения метрики.
type GetMetricService struct {
	getter Getter[string, *domain.Metrics] // используется для чтения метрик
}

// NewGetMetricService конструктор для GetMetricService.
func NewGetMetricService(getter Getter[string, *domain.Metrics]) *GetMetricService {
	return &GetMetricService{
		getter: getter,
	}
}

// GetMetric получает метрику по ID и типу.
func (s *GetMetricService) GetMetric(id string, mtype domain.MetricType) (*domain.Metrics, error) {
	key := id + ":" + string(mtype)
	metric, exists := s.getter.Get(key)
	if !exists {
		return nil, ErrMetricNotFound
	}
	return metric, nil
}
