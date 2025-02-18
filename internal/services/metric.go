package services

import (
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/types"

	"go.uber.org/zap"
)

// Saver управляет операцией записи в хранилище.
type Saver[K comparable, V any] interface {
	Save(key K, value V) bool
}

// Getter управляет операцией чтения из хранилища.
type Getter[K comparable, V any] interface {
	Get(key K) (V, bool)
}

type UpdateMetricsService struct {
	gaugeSaver    Saver[string, float64]
	gaugeGetter   Getter[string, float64]
	counterSaver  Saver[string, int64]
	counterGetter Getter[string, int64]
}

// NewUpdateMetricsService создает новый сервис для работы с метриками.
func NewUpdateMetricsService(
	gaugeSaver Saver[string, float64],
	gaugeGetter Getter[string, float64],
	counterSaver Saver[string, int64],
	counterGetter Getter[string, int64],
) *UpdateMetricsService {
	logger.Logger.Info("UpdateMetricsService initialized")

	return &UpdateMetricsService{
		gaugeSaver:    gaugeSaver,
		gaugeGetter:   gaugeGetter,
		counterSaver:  counterSaver,
		counterGetter: counterGetter,
	}
}

// Update обновляет метрику в зависимости от типа и переданных данных.
func (s *UpdateMetricsService) Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error) {
	logger.Logger.Debug("Received metric update request",
		zap.String("metric_id", req.ID),
		zap.String("metric_type", string(req.MType)))

	switch req.MType {
	case types.Gauge:
		s.gaugeSaver.Save(req.ID, *req.Value)
		logger.Logger.Info("Gauge metric updated",
			zap.String("metric_id", req.ID),
			zap.Float64("value", *req.Value))

		return &types.UpdateMetricsResponse{
			UpdateMetricsRequest: *req,
		}, nil

	case types.Counter:
		existingValue, exists := s.counterGetter.Get(req.ID)
		if !exists {
			existingValue = int64(0)
			logger.Logger.Debug("Counter metric not found, initializing to 0",
				zap.String("metric_id", req.ID))
		}

		newValue := existingValue + *req.Delta
		s.counterSaver.Save(req.ID, newValue)

		logger.Logger.Info("Counter metric updated",
			zap.String("metric_id", req.ID),
			zap.Int64("previous_value", existingValue),
			zap.Int64("delta", *req.Delta),
			zap.Int64("new_value", newValue))

		return &types.UpdateMetricsResponse{
			UpdateMetricsRequest: *req,
		}, nil

	default:
		logger.Logger.Warn("Received unknown metric type",
			zap.String("metric_id", req.ID),
			zap.String("metric_type", string(req.MType)))

		return nil, nil
	}
}
