package services

import (
	"errors"
	"go-metrics-alerting/internal/domain"
)

var ErrUnknownMetricType = errors.New("unknown metric type")

// UpdateMetricStrategy интерфейс для стратегий обновления метрик.
type UpdateMetricStrategy interface {
	UpdateMetric(metric *domain.Metrics) *domain.Metrics
}

type UpdateMetricService struct {
	updateCounterStrategy UpdateMetricStrategy
	updateGaugeStrategy   UpdateMetricStrategy
}

// NewUpdateMetricService создает новый сервис для обновления метрик.
func NewUpdateMetricService(
	updateCounterStrategy UpdateMetricStrategy,
	updateGaugeStrategy UpdateMetricStrategy,
) *UpdateMetricService {
	return &UpdateMetricService{
		updateCounterStrategy: updateCounterStrategy,
		updateGaugeStrategy:   updateGaugeStrategy,
	}
}

// UpdateMetric обновляет метрику, выбирая соответствующую стратегию в зависимости от типа метрики.
func (s *UpdateMetricService) UpdateMetric(metric *domain.Metrics) (*domain.Metrics, error) {
	switch metric.MType {
	case domain.Counter:
		return s.updateCounterStrategy.UpdateMetric(metric), nil
	case domain.Gauge:
		return s.updateGaugeStrategy.UpdateMetric(metric), nil
	default:
		return nil, ErrUnknownMetricType
	}
}
