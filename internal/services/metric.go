package services

import (
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/types"
	"strconv"

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

// Getter управляет операцией чтения из хранилища.
type Ranger[K comparable, V any] interface {
	Range(callback func(key K, value V) bool)
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

// GetMetricValueServiceImpl структура для реализации сервиса получения значения метрики.
type GetMetricValueService struct {
	gaugeGetter   Getter[string, float64]
	counterGetter Getter[string, int64]
}

// NewGetMetricValueService создаёт новый экземпляр сервиса получения метрик.
func NewGetMetricValueService(
	gaugeGetter Getter[string, float64],
	counterGetter Getter[string, int64],
) *GetMetricValueService {
	logger.Logger.Info("GetMetricValueService initialized")
	return &GetMetricValueService{
		gaugeGetter:   gaugeGetter,
		counterGetter: counterGetter,
	}
}

// GetMetricValue получает значение метрики по её ID и возвращает результат с нужным типом.
func (s *GetMetricValueService) GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error) {
	logger.Logger.Debug("Received request to get metric value",
		zap.String("metric_id", req.ID),
		zap.String("metric_type", string(req.MType)))

	switch req.MType {
	case types.Gauge:
		if value, exists := s.gaugeGetter.Get(req.ID); exists {
			return &types.GetMetricValueResponse{
				ID:    req.ID,
				Value: strconv.FormatFloat(value, 'f', -1, 64),
			}, nil
		}
	case types.Counter:
		if value, exists := s.counterGetter.Get(req.ID); exists {
			return &types.GetMetricValueResponse{
				ID:    req.ID,
				Value: strconv.FormatInt(value, 10),
			}, nil
		}
	default:
		logger.Logger.Warn("Unknown metric type received",
			zap.String("metric_id", req.ID),
			zap.String("metric_type", string(req.MType)))
		return nil, nil
	}

	logger.Logger.Warn("Metric not found",
		zap.String("metric_id", req.ID),
		zap.String("metric_type", string(req.MType)))
	return nil, nil
}

// GetAllMetricValuesServiceImpl структура для реализации сервиса получения всех метрик.
type GetAllMetricValuesService struct {
	gaugeRanger   Ranger[string, float64]
	counterRanger Ranger[string, int64]
}

// NewGetAllMetricValuesService создаёт новый экземпляр сервиса получения всех метрик.
func NewGetAllMetricValuesService(
	gaugeRanger Ranger[string, float64],
	counterRanger Ranger[string, int64],
) *GetAllMetricValuesService {
	logger.Logger.Info("GetAllMetricValuesService initialized")
	return &GetAllMetricValuesService{
		gaugeRanger:   gaugeRanger,
		counterRanger: counterRanger,
	}
}

// GetAllMetricValues получает все значения метрик.
func (s *GetAllMetricValuesService) GetAllMetricValues() []*types.GetMetricValueResponse {
	logger.Logger.Debug("Received request to get all metric values")

	var metrics []*types.GetMetricValueResponse

	// Считываем все гейджи
	s.gaugeRanger.Range(func(key string, value float64) bool {
		metrics = append(metrics, &types.GetMetricValueResponse{
			ID:    key,
			Value: strconv.FormatFloat(value, 'f', -1, 64),
		})
		return true
	})

	// Считываем все счётчики
	s.counterRanger.Range(func(key string, value int64) bool {
		metrics = append(metrics, &types.GetMetricValueResponse{
			ID:    key,
			Value: strconv.FormatInt(value, 10),
		})
		return true
	})

	// Возвращаем все метрики
	return metrics
}
