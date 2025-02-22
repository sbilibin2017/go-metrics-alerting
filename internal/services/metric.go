package services

import (
	"go-metrics-alerting/internal/domain"
)

// Saver определяет методы для сохранения данных в хранилище.
type Saver interface {
	Save(key string, value *domain.Metrics)
}

// KeyEncoder интерфейс для кодирования ключей.
type UpdateMetricStrategy interface {
	Update(metric *domain.Metrics) *domain.Metrics
}

// UpdateMetricService фасад для двух сервисов
type UpdateMetricService struct {
	updateCounterStrategy UpdateMetricStrategy
	updateGaugeStrategy   UpdateMetricStrategy
}

// NewUpdateMetricService создаёт новый фасадный сервис для работы с метриками
func NewUpdateMetricService(
	updateCounterStrategy UpdateMetricStrategy,
	updateGaugeStrategy UpdateMetricStrategy,
) *UpdateMetricService {
	return &UpdateMetricService{
		updateCounterStrategy: updateCounterStrategy,
		updateGaugeStrategy:   updateGaugeStrategy,
	}
}

// Update обновляет значение метрики в зависимости от её типа
func (s *UpdateMetricService) UpdateMetric(metric *domain.Metrics) *domain.Metrics {
	// Определение стратегии по типу метрики
	switch metric.MType {
	case domain.Counter:
		return s.updateCounterStrategy.Update(metric)
	case domain.Gauge:
		return s.updateGaugeStrategy.Update(metric)
	default:
		return nil
	}

}

// Getter определяет методы для получения данных из хранилища.
type Getter interface {
	Get(key string) *domain.Metrics
}

// KeyEncoder интерфейс для кодирования ключей.
type KeyEncoder interface {
	Encode(id, mtype string) string
}

// GetMetricService для получения одной метрики по её ID
type GetMetricService struct {
	getter     Getter
	keyEncoder KeyEncoder
}

// NewGetMetricService создаёт новый сервис для получения метрики по ID
func NewGetMetricService(getter Getter, keyEncoder KeyEncoder) *GetMetricService {
	return &GetMetricService{getter: getter, keyEncoder: keyEncoder}
}

// Get возвращает метрику по её ID
func (s *GetMetricService) Get(id string, mtype domain.MType) *domain.Metrics {
	key := s.keyEncoder.Encode(id, string(mtype))
	return s.getter.Get(key)
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger interface {
	Range(callback func(key string, value *domain.Metrics) bool)
}

// GetAllMetricsService для получения всех метрик с использованием Ranger
type GetAllMetricsService struct {
	ranger Ranger
}

// NewGetAllMetricsService создаёт новый сервис для получения всех метрик
func NewGetAllMetricsService(ranger Ranger) *GetAllMetricsService {
	return &GetAllMetricsService{ranger: ranger}
}

// GetAll перебирает все метрики и возвращает их как срез.
func (s *GetAllMetricsService) GetAll() []*domain.Metrics {
	var result []*domain.Metrics
	s.ranger.Range(func(key string, value *domain.Metrics) bool {
		result = append(result, value)
		return true
	})
	return result
}
