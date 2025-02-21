package strategies

import "go-metrics-alerting/internal/domain"

// Saver определяет методы для сохранения данных в хранилище.
type Saver interface {
	Save(key string, value *domain.Metrics)
}

// Getter определяет методы для получения данных из хранилища.
type Getter interface {
	Get(key string) *domain.Metrics
}

// KeyEncoder интерфейс для кодирования ключей.
type KeyEncoder interface {
	Encode(id, mtype string) string
}

// UpdateGaugeMetricService для обновления метрик типа gauge
type UpdateGaugeMetricStrategy struct {
	saver      Saver
	getter     Getter
	keyEncoder KeyEncoder
}

// NewUpdateGaugeMetricService создаёт новый сервис для работы с gauge метриками
func NewUpdateGaugeMetricStrategy(saver Saver, getter Getter, keyEncoder KeyEncoder) *UpdateGaugeMetricStrategy {
	return &UpdateGaugeMetricStrategy{saver: saver, getter: getter, keyEncoder: keyEncoder}
}

// Update обновляет значение метрики типа gauge (замещает старое значение)
func (s *UpdateGaugeMetricStrategy) Update(metric *domain.Metrics) *domain.Metrics {
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
type UpdateCounterMetricStrategy struct {
	saver      Saver
	getter     Getter
	keyEncoder KeyEncoder
}

// NewUpdateCounterMetricService создаёт новый сервис для работы с counter метриками
func NewUpdateCounterMetricStrategy(saver Saver, getter Getter, keyEncoder KeyEncoder) *UpdateCounterMetricStrategy {
	return &UpdateCounterMetricStrategy{saver: saver, getter: getter, keyEncoder: keyEncoder}
}

// Update обновляет значение метрики типа counter (прибавляет к текущему значению)
func (s *UpdateCounterMetricStrategy) Update(metric *domain.Metrics) *domain.Metrics {
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
