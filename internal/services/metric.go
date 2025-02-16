package services

import (
	"errors"
	"go-metrics-alerting/internal/types"
	"strconv"
)

var (
	ErrMetricIsNotUpdated error = errors.New("metric is not updated")
	ErrMetricIsNotFound   error = errors.New("metric is not found")
)

type Setter interface {
	Set(key string, value string) error
}

type Getter interface {
	Get(key string) (string, error)
}

type Ranger interface {
	Range(callback func(key string, value string) bool)
}

// UpdateMetricService сервис для обновления метрик.
type UpdateMetricService struct {
	stringGetter Getter
	stringSetter Setter
}

// NewUpdateMetricService создаёт новый сервис для обновления метрик.
func NewUpdateMetricService(
	stringGetter Getter,
	stringSetter Setter,
) *UpdateMetricService {
	return &UpdateMetricService{
		stringGetter: stringGetter,
		stringSetter: stringSetter,
	}
}

// UpdateMetric обновляет или сохраняет метрику в зависимости от её типа.
func (s *UpdateMetricService) UpdateMetric(req *types.MetricsRequest) (*types.MetricsResponse, error) {
	currentValue, err := s.stringGetter.Get(req.ID)
	if err != nil {
		currentValue = "0"
	}

	switch req.MType {
	case types.Counter:
		intValue, _ := strconv.ParseInt(currentValue, 10, 64)
		*req.Delta += intValue
		err := s.stringSetter.Set(req.ID, strconv.FormatInt(*req.Delta, 10))
		if err != nil {
			return nil, ErrMetricIsNotUpdated
		}
	case types.Gauge:
		err := s.stringSetter.Set(req.ID, strconv.FormatFloat(*req.Value, 'f', -1, 64))
		if err != nil {
			return nil, ErrMetricIsNotUpdated
		}
	}

	return &types.MetricsResponse{MetricsRequest: *req}, nil
}

// GetMetricService сервис для получения метрик.
type GetMetricService struct {
	stringGetter Getter
}

// NewGetMetricService создаёт новый сервис для получения метрик.
func NewGetMetricService(
	stringGetter Getter,
) *GetMetricService {
	return &GetMetricService{
		stringGetter: stringGetter,
	}
}

// GetMetric получает метрику по её ID.
func (s *GetMetricService) GetMetric(metric *types.MetricsRequest) (*string, error) {
	value, err := s.stringGetter.Get(metric.ID)
	if err != nil {
		return nil, ErrMetricIsNotFound
	}
	return &value, nil
}

// ListMetricsService сервис для получения списка метрик.
type ListMetricsService struct {
	stringRanger Ranger
}

// NewListMetricsService создаёт новый сервис для получения списка метрик.
func NewListMetricsService(
	stringRanger Ranger,
) *ListMetricsService {
	return &ListMetricsService{
		stringRanger: stringRanger,
	}
}

// ListMetrics возвращает список всех метрик.
func (s *ListMetricsService) ListMetrics() []*types.MetricValueResponse {
	var metrics []*types.MetricValueResponse
	s.stringRanger.Range(func(key string, value string) bool {
		metrics = append(metrics, &types.MetricValueResponse{
			ID:    key,
			Value: value,
		})
		return true
	})
	return metrics
}
