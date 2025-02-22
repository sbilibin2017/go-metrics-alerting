package services

import (
	"go-metrics-alerting/internal/domain"

	"go.uber.org/zap"
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
	logger                *zap.Logger
}

// NewUpdateMetricService создаёт новый фасадный сервис для работы с метриками
func NewUpdateMetricService(
	updateCounterStrategy UpdateMetricStrategy,
	updateGaugeStrategy UpdateMetricStrategy,
	logger *zap.Logger, // добавляем логгер
) *UpdateMetricService {
	return &UpdateMetricService{
		updateCounterStrategy: updateCounterStrategy,
		updateGaugeStrategy:   updateGaugeStrategy,
		logger:                logger,
	}
}

// Update обновляет значение метрики в зависимости от её типа
func (s *UpdateMetricService) UpdateMetric(metric *domain.Metrics) *domain.Metrics {
	// Логируем получение метрики
	s.logger.Info("Received request to update metric", zap.String("ID", metric.ID), zap.String("Type", string(metric.MType)))

	// Определение стратегии по типу метрики
	switch metric.MType {
	case domain.Counter:
		s.logger.Info("Using Counter strategy for update", zap.String("ID", metric.ID))
		return s.updateCounterStrategy.Update(metric)
	case domain.Gauge:
		s.logger.Info("Using Gauge strategy for update", zap.String("ID", metric.ID))
		return s.updateGaugeStrategy.Update(metric)
	default:
		s.logger.Warn("Unsupported metric type", zap.String("Type", string(metric.MType)))
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
	logger     *zap.Logger
}

// NewGetMetricService создаёт новый сервис для получения метрики по ID
func NewGetMetricService(getter Getter, keyEncoder KeyEncoder, logger *zap.Logger) *GetMetricService {
	return &GetMetricService{getter: getter, keyEncoder: keyEncoder, logger: logger}
}

// Get возвращает метрику по её ID
func (s *GetMetricService) GetMetric(id string, mtype domain.MType) *domain.Metrics {
	key := s.keyEncoder.Encode(id, string(mtype))
	s.logger.Info("Fetching metric", zap.String("ID", id), zap.String("Type", string(mtype)))
	return s.getter.Get(key)
}

// Ranger управляет операцией перебора элементов в хранилище.
type Ranger interface {
	Range(callback func(key string, value *domain.Metrics) bool)
}

// GetAllMetricsService для получения всех метрик с использованием Ranger
type GetAllMetricsService struct {
	ranger Ranger
	logger *zap.Logger
}

// NewGetAllMetricsService создаёт новый сервис для получения всех метрик
func NewGetAllMetricsService(ranger Ranger, logger *zap.Logger) *GetAllMetricsService {
	return &GetAllMetricsService{ranger: ranger, logger: logger}
}

// GetAll перебирает все метрики и возвращает их как срез.
func (s *GetAllMetricsService) GetAllMetrics() []*domain.Metrics {
	var result []*domain.Metrics
	s.logger.Info("Fetching all metrics")

	s.ranger.Range(func(key string, value *domain.Metrics) bool {
		s.logger.Debug("Adding metric to result", zap.String("ID", value.ID), zap.String("Type", string(value.MType)))
		result = append(result, value)
		return true
	})

	s.logger.Info("Fetched all metrics", zap.Int("Count", len(result)))
	return result
}
