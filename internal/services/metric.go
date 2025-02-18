package services

import (
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/types"
	"strconv"

	"go.uber.org/zap"
)

type Saver interface {
	Save(key, value string) bool
}

type Getter interface {
	Get(key string) (string, bool)
}

type Ranger interface {
	Range(callback func(key, value string) bool)
}

type UpdateMetricsService struct {
	saver  Saver
	getter Getter
}

func NewUpdateMetricsService(saver Saver, getter Getter) *UpdateMetricsService {
	logger.Logger.Info("UpdateMetricsService initialized")
	return &UpdateMetricsService{
		saver:  saver,
		getter: getter,
	}
}

func (s *UpdateMetricsService) Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error) {
	logger.Logger.Debug("Received metric update request",
		zap.String("metric_id", req.ID),
		zap.String("metric_type", string(req.MType)))

	switch req.MType {
	case types.Gauge:
		s.saver.Save(req.ID, strconv.FormatFloat(*req.Value, 'f', -1, 64))

	case types.Counter:
		existingValue := int64(0)
		if existingValueStr, exists := s.getter.Get(req.ID); exists {
			existingValue, _ = strconv.ParseInt(existingValueStr, 10, 64)
		}
		*req.Delta += existingValue
		s.saver.Save(req.ID, strconv.FormatInt(*req.Delta, 10))

	default:
		return nil, nil
	}

	return &types.UpdateMetricsResponse{UpdateMetricsRequest: *req}, nil
}

type GetMetricValueService struct {
	getter Getter
}

func NewGetMetricValueService(getter Getter) *GetMetricValueService {
	logger.Logger.Info("GetMetricValueService initialized")
	return &GetMetricValueService{
		getter: getter,
	}
}

func (s *GetMetricValueService) GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error) {
	logger.Logger.Debug("Received request to get metric value",
		zap.String("metric_id", req.ID),
		zap.String("metric_type", string(req.MType)))

	if valueStr, exists := s.getter.Get(req.ID); exists {
		return &types.GetMetricValueResponse{ID: req.ID, Value: valueStr}, nil
	}

	return nil, nil
}

type GetAllMetricValuesService struct {
	ranger Ranger
}

func NewGetAllMetricValuesService(ranger Ranger) *GetAllMetricValuesService {
	logger.Logger.Info("GetAllMetricValuesService initialized")
	return &GetAllMetricValuesService{ranger: ranger}
}

func (s *GetAllMetricValuesService) GetAllMetricValues() []*types.GetMetricValueResponse {
	logger.Logger.Debug("Received request to get all metric values")
	var metrics []*types.GetMetricValueResponse

	s.ranger.Range(func(key, value string) bool {
		metrics = append(metrics, &types.GetMetricValueResponse{ID: key, Value: value})
		return true
	})

	return metrics
}
