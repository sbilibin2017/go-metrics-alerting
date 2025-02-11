package services

import (
	"context"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"strconv"
)

// MetricRepository — интерфейс хранилища метрик.
type MetricSaverRepository interface {
	Save(ctx context.Context, metricType, metricName, value string) error
}

// MetricRepository — интерфейс хранилища метрик.
type MetricGetterRepository interface {
	Get(ctx context.Context, metricType, metricName string) (string, error)
}

// MetricRepository — интерфейс хранилища метрик.
type MetricListerRepository interface {
	GetAll(ctx context.Context) [][]string
}

// MetricRepository — интерфейс хранилища метрик.
type MetricRepository interface {
	MetricSaverRepository
	MetricGetterRepository
	MetricListerRepository
}

const (
	MetricEmptyString  string = ""
	DefaultMetricValue string = "0"
)

// MetricService — сервис для работы с метриками.
type MetricService struct {
	MetricRepository MetricRepository
}

// UpdateMetricValue обновляет значение метрики.
func (s *MetricService) UpdateMetric(
	ctx context.Context,
	req *types.UpdateMetricValueRequest,
) error {
	// Проверка на отсутствие имени метрики
	if req.Name == MetricEmptyString {
		return &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: "metric name is required",
		}
	}

	// Проверка на отсутствие имени метрики
	if string(req.Type) == MetricEmptyString {
		return &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: "metric type is required",
		}
	}

	// Проверка на корректность типа метрики
	if req.Type != types.Counter && req.Type != types.Gauge {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "invalid metric type",
		}
	}

	currentValue, err := s.MetricRepository.Get(ctx, string(req.Type), req.Name)
	if err != nil {
		currentValue = DefaultMetricValue
	}

	var value string
	switch req.Type {
	case types.Counter:
		newVal, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			return &apierror.APIError{
				Code:    http.StatusBadRequest,
				Message: "invalid counter value",
			}
		}
		curVal, _ := strconv.ParseInt(currentValue, 10, 64)
		newVal += curVal
		value = strconv.FormatInt(newVal, 10)

	case types.Gauge:
		newVal, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return &apierror.APIError{
				Code:    http.StatusBadRequest,
				Message: "invalid gauge value",
			}
		}
		value = strconv.FormatFloat(newVal, 'f', -1, 64)
	}

	err = s.MetricRepository.Save(ctx, string(req.Type), req.Name, value)
	if err != nil {
		return &apierror.APIError{
			Code:    http.StatusInternalServerError,
			Message: "value is not saved",
		}
	}
	return nil
}

// GetMetricValue возвращает значение метрики по имени и типу.
func (s *MetricService) GetMetric(ctx context.Context, req *types.GetMetricValueRequest) (string, error) {
	currentValue, err := s.MetricRepository.Get(ctx, string(req.Type), req.Name)
	if err != nil {
		return MetricEmptyString, &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: "value not found",
		}
	}
	return currentValue, nil
}

// GetAllMetrics возвращает список всех метрик.
func (s *MetricService) ListMetrics(ctx context.Context) []*types.MetricResponse {
	metricsList := s.MetricRepository.GetAll(ctx)
	var metrics []*types.MetricResponse

	if len(metricsList) == 0 {
		return []*types.MetricResponse{}
	}

	for _, metric := range metricsList {
		metrics = append(metrics, &types.MetricResponse{
			Name:  metric[1],
			Value: metric[2],
		})
	}

	return metrics
}
