package services

import (
	"errors"
	"go-metrics-alerting/internal/domain"
)

// Ошибка для неудачного обновления метрики
var ErrUpdateFailed = errors.New("failed to update metric")

// UpdateMetricStrategy интерфейс для стратегии обновления метрики
type UpdateMetricStrategy interface {
	Update(metric *domain.Metric) (*domain.Metric, bool)
}

// UpdateMetricsService структура для обработки обновлений метрик
type UpdateMetricsService struct {
	strategy UpdateMetricStrategy
}

// UpdateMetricValue обновляет значение метрики, используя стратегию
func (s *UpdateMetricsService) UpdateMetricValue(metric *domain.Metric) (*domain.Metric, error) {
	updatedMetric, ok := s.strategy.Update(metric)
	if !ok {
		return nil, ErrUpdateFailed
	}
	return updatedMetric, nil
}

// Ошибка для случая, когда метрика не найдена
var ErrMetricNotFound = errors.New("metric not found")

// MetricGetter интерфейс для получения метрик.
type Getter interface {
	Get(id string) (string, bool)
}

type KeyEncoder interface {
	Encode(id string, mtype string) string
}

// Сервис для получения значения метрики по ID
type GetMetricValueService struct {
	getter  Getter
	encoder KeyEncoder
}

// Метод для получения значения метрики
func (s *GetMetricValueService) GetMetricValue(id string, mType domain.MType) (*domain.Metric, error) {
	if valueStr, exists := s.getter.Get(s.encoder.Encode(id, string(mType))); exists {
		return &domain.Metric{
			ID:    id,
			MType: mType,
			Value: valueStr,
		}, nil
	}
	return nil, ErrMetricNotFound
}

// MetricGetter интерфейс для получения метрик.
type Ranger interface {
	Range(callback func(id, value string) bool)
}
type KeyDecoder interface {
	Decode(key string) (id string, mtype string, ok bool)
}

// Сервис для получения всех метрик
type GetAllMetricValuesService struct {
	ranger  Ranger
	decoder KeyDecoder
}

// Метод для получения всех метрик
func (s *GetAllMetricValuesService) GetAllMetricValues() []*domain.Metric {
	var metrics []*domain.Metric
	s.ranger.Range(func(id, value string) bool {
		id, mType, ok := s.decoder.Decode(id)
		if !ok {
			return true
		}
		metrics = append(metrics, &domain.Metric{
			ID:    id,
			MType: domain.MType(mType),
			Value: value,
		})

		return true
	})
	return metrics
}
