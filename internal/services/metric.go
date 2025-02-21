package services

import (
	"go-metrics-alerting/internal/types"
)

// Saver определяет методы для сохранения данных в хранилище.
type Saver interface {
	Save(key string, value *types.Metrics)
}

// Getter определяет методы для получения данных из хранилища.
type Getter interface {
	Get(key string) *types.Metrics
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger interface {
	Range(callback func(key string, value *types.Metrics) bool)
}

// KeyEncoder интерфейс для кодирования ключей.
type KeyEncoder interface {
	Encode(id, mtype string) string
}

// UpdateGaugeMetricService для обновления метрик типа gauge
type UpdateGaugeMetricService struct {
	saver      Saver
	getter     Getter
	keyEncoder KeyEncoder
}

// NewUpdateGaugeMetricService создаёт новый сервис для работы с gauge метриками
func NewUpdateGaugeMetricService(saver Saver, getter Getter, keyEncoder KeyEncoder) *UpdateGaugeMetricService {
	return &UpdateGaugeMetricService{saver: saver, getter: getter, keyEncoder: keyEncoder}
}

// Update обновляет значение метрики типа gauge (замещает старое значение)
func (s *UpdateGaugeMetricService) Update(metric *types.Metrics) *types.Metrics {
	key := s.keyEncoder.Encode(metric.ID, string(metric.MType))
	existingMetric := s.getter.Get(key)
	if existingMetric == nil {
		s.saver.Save(key, metric)
		return metric
	}
	existingMetric.Value = metric.Value
	s.saver.Save(key, existingMetric)
	return existingMetric
}

// UpdateCounterMetricService для обновления метрик типа counter
type UpdateCounterMetricService struct {
	saver      Saver
	getter     Getter
	keyEncoder KeyEncoder
}

// NewUpdateCounterMetricService создаёт новый сервис для работы с counter метриками
func NewUpdateCounterMetricService(saver Saver, getter Getter, keyEncoder KeyEncoder) *UpdateCounterMetricService {
	return &UpdateCounterMetricService{saver: saver, getter: getter, keyEncoder: keyEncoder}
}

// Update обновляет значение метрики типа counter (прибавляет к текущему значению)
func (s *UpdateCounterMetricService) Update(metric *types.Metrics) *types.Metrics {
	key := s.keyEncoder.Encode(metric.ID, string(metric.MType))
	existingMetric := s.getter.Get(key)
	if existingMetric == nil {
		s.saver.Save(key, metric)
		return metric
	}
	*existingMetric.Delta += *metric.Delta
	s.saver.Save(key, existingMetric)
	return existingMetric
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
func (s *UpdateMetricService) Update(metric *types.Metrics) *types.Metrics {
	switch metric.MType {
	case types.Counter:
		return s.counterService.Update(metric)
	case types.Gauge:
		return s.gaugeService.Update(metric)
	default:
		return nil
	}
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
func (s *GetMetricService) Get(id string, mtype types.MType) *types.Metrics {
	key := s.keyEncoder.Encode(id, string(mtype))
	return s.getter.Get(key)
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
func (s *GetAllMetricsService) GetAll() []*types.Metrics {
	var result []*types.Metrics
	s.ranger.Range(func(key string, value *types.Metrics) bool {
		result = append(result, value)
		return true
	})
	return result
}
