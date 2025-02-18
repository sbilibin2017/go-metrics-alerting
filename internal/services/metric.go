package services

import (
	"go-metrics-alerting/internal/api/types"
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
	return &UpdateMetricsService{
		gaugeSaver:    gaugeSaver,
		gaugeGetter:   gaugeGetter,
		counterSaver:  counterSaver,
		counterGetter: counterGetter,
	}
}

// Update обновляет метрику в зависимости от типа и переданных данных.
func (s *UpdateMetricsService) Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error) {
	switch req.MType {
	case types.Gauge:
		s.gaugeSaver.Save(req.ID, *req.Value)
		return &types.UpdateMetricsResponse{
			UpdateMetricsRequest: *req,
		}, nil

	case types.Counter:
		existingValue, exists := s.counterGetter.Get(req.ID)
		if !exists {
			existingValue = int64(0)
		}
		s.counterSaver.Save(req.ID, existingValue+*req.Delta)
		return &types.UpdateMetricsResponse{
			UpdateMetricsRequest: *req,
		}, nil

	default:
		return nil, nil
	}
}
