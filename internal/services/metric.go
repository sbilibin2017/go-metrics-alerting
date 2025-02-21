package services

import (
	"go-metrics-alerting/internal/types"
)

// Saver определяет методы для сохранения данных в хранилище.
type Saver interface {
	Save(key string, value *types.Metrics) bool
}

// Getter определяет методы для получения данных из хранилища.
type Getter interface {
	Get(key string) (*types.Metrics, bool)
}

// UpdateGaugeMetricService для обновления метрик типа gauge
type UpdateGaugeMetricService struct {
	saver  Saver
	getter Getter
}

// NewUpdateGaugeMetricService создаёт новый сервис для работы с gauge метриками
func NewUpdateGaugeMetricService(saver Saver, getter Getter) *UpdateGaugeMetricService {
	return &UpdateGaugeMetricService{saver: saver, getter: getter}
}

// Update обновляет значение метрики типа gauge (замещает старое значение)
func (s *UpdateGaugeMetricService) Update(metric *types.Metrics) (*types.Metrics, bool) {
	existingMetric, exists := s.getter.Get(metric.ID)
	if !exists {
		s.saver.Save(metric.ID, metric)
		return metric, true
	}
	existingMetric.Value = metric.Value
	s.saver.Save(existingMetric.ID, existingMetric)
	return existingMetric, true
}

// UpdateCounterMetricService для обновления метрик типа counter
type UpdateCounterMetricService struct {
	saver  Saver
	getter Getter
}

// NewUpdateCounterMetricService создаёт новый сервис для работы с counter метриками
func NewUpdateCounterMetricService(saver Saver, getter Getter) *UpdateCounterMetricService {
	return &UpdateCounterMetricService{saver: saver, getter: getter}
}

// Update обновляет значение метрики типа counter (прибавляет к текущему значению)
func (s *UpdateCounterMetricService) Update(metric *types.Metrics) (*types.Metrics, bool) {
	existingMetric, exists := s.getter.Get(metric.ID)
	if !exists {
		s.saver.Save(metric.ID, metric)
		return metric, true
	}
	*existingMetric.Delta += *metric.Delta
	s.saver.Save(existingMetric.ID, existingMetric)
	return existingMetric, true
}

// UpdateMetricService фасад для двух сервисов
type UpdateMetricService struct {
	gaugeService   *UpdateGaugeMetricService
	counterService *UpdateCounterMetricService
}

// NewUpdateMetricService создаёт новый фасадный сервис для работы с метриками
func NewUpdateMetricService(gaugeService *UpdateGaugeMetricService, counterService *UpdateCounterMetricService) *UpdateMetricService {
	return &UpdateMetricService{gaugeService: gaugeService, counterService: counterService}
}

// Update обновляет значение метрики в зависимости от её типа
func (s *UpdateMetricService) Update(metric *types.Metrics) (*types.Metrics, bool) {
	switch metric.MType {
	case types.Counter:
		return s.counterService.Update(metric)
	case types.Gauge:
		return s.gaugeService.Update(metric)
	default:
		return nil, false
	}
}
