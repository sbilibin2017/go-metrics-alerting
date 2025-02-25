package strategies

import (
	"go-metrics-alerting/internal/domain"
)

// Getter интерфейс для чтения из обобщённого хранилища.
type Getter[K comparable, V any] interface {
	Get(key K) (V, bool)
}

// Saver интерфейс для записи данных в обобщённое хранилище.
type Saver[K comparable, V any] interface {
	Save(key K, value V)
}

type UpdateGaugeMetricStrategy struct {
	saver  Saver[string, *domain.Metrics]
	getter Getter[string, *domain.Metrics]
}

func NewUpdateGaugeMetricStrategy(
	saver Saver[string, *domain.Metrics],
	getter Getter[string, *domain.Metrics],
) *UpdateGaugeMetricStrategy {
	return &UpdateGaugeMetricStrategy{saver: saver, getter: getter}
}

func (s *UpdateGaugeMetricStrategy) UpdateMetric(metric *domain.Metrics) *domain.Metrics {
	key := metric.ID + ":" + string(metric.MType)
	s.saver.Save(key, metric)
	return metric
}

type UpdateCounterMetricStrategy struct {
	saver  Saver[string, *domain.Metrics]
	getter Getter[string, *domain.Metrics]
}

func NewUpdateCounterMetricStrategy(
	saver Saver[string, *domain.Metrics],
	getter Getter[string, *domain.Metrics],
) *UpdateCounterMetricStrategy {
	return &UpdateCounterMetricStrategy{saver: saver, getter: getter}
}

func (s *UpdateCounterMetricStrategy) UpdateMetric(metric *domain.Metrics) *domain.Metrics {
	key := metric.ID + ":" + string(metric.MType)
	existingMetric, exists := s.getter.Get(key)

	if exists {
		*existingMetric.Delta += *metric.Delta
		s.saver.Save(key, existingMetric)
		return existingMetric
	} else {
		s.saver.Save(key, metric)
		return metric
	}
}
