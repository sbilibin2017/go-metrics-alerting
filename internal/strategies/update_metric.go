package strategies

import (
	"go-metrics-alerting/internal/domain"

	"go.uber.org/zap"
)

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
	logger     *zap.Logger
}

// NewUpdateGaugeMetricStrategy создаёт новый сервис для работы с gauge метриками
func NewUpdateGaugeMetricStrategy(saver Saver, getter Getter, keyEncoder KeyEncoder, logger *zap.Logger) *UpdateGaugeMetricStrategy {
	return &UpdateGaugeMetricStrategy{saver: saver, getter: getter, keyEncoder: keyEncoder, logger: logger}
}

// Update обновляет значение метрики типа gauge (замещает старое значение)
func (s *UpdateGaugeMetricStrategy) Update(metric *domain.Metrics) *domain.Metrics {
	key := s.keyEncoder.Encode(metric.ID, string(metric.MType))
	existingMetric := s.getter.Get(key)
	if existingMetric == nil {
		s.logger.Info("Saving new gauge metric", zap.String("metricID", metric.ID))
		s.saver.Save(key, metric)
		return metric
	}
	// Логирование обновления существующей метрики
	s.logger.Info("Updating existing gauge metric", zap.String("metricID", metric.ID), zap.Float64("oldValue", *existingMetric.Value), zap.Float64("newValue", *metric.Value))
	existingMetric.Value = metric.Value
	s.saver.Save(key, existingMetric)
	return existingMetric
}

// UpdateCounterMetricService для обновления метрик типа counter
type UpdateCounterMetricStrategy struct {
	saver      Saver
	getter     Getter
	keyEncoder KeyEncoder
	logger     *zap.Logger
}

// NewUpdateCounterMetricStrategy создаёт новый сервис для работы с counter метриками
func NewUpdateCounterMetricStrategy(saver Saver, getter Getter, keyEncoder KeyEncoder, logger *zap.Logger) *UpdateCounterMetricStrategy {
	return &UpdateCounterMetricStrategy{saver: saver, getter: getter, keyEncoder: keyEncoder, logger: logger}
}

// Update обновляет значение метрики типа counter (прибавляет к текущему значению)
func (s *UpdateCounterMetricStrategy) Update(metric *domain.Metrics) *domain.Metrics {
	key := s.keyEncoder.Encode(metric.ID, string(metric.MType))
	existingMetric := s.getter.Get(key)
	if existingMetric == nil {
		s.logger.Info("Saving new counter metric", zap.String("metricID", metric.ID))
		s.saver.Save(key, metric)
		return metric
	}
	// Логирование обновления счетчика
	s.logger.Info("Updating existing counter metric", zap.String("metricID", metric.ID), zap.Int64("oldDelta", *existingMetric.Delta), zap.Int64("newDelta", *metric.Delta))
	*existingMetric.Delta += *metric.Delta
	s.saver.Save(key, existingMetric)
	return existingMetric
}
