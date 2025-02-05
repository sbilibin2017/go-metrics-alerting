package services

import (
	"go-metrics-alerting/internal/metric/handlers"
	"go-metrics-alerting/pkg/apierror"

	"net/http"
)

const (
	getValueEmptyString string = ""
)

// Интерфейс для получения значения метрики
type MetricStorageGetter interface {
	Get(metricType string, metricName string) (string, error)
}

// Сервис для получения значения метрики
type GetMetricValueService struct {
	metricRepository MetricStorageGetter
}

// Новый сервис для получения значения метрики
func NewGetMetricValueService(metricRepository MetricStorageGetter) *GetMetricValueService {
	return &GetMetricValueService{
		metricRepository: metricRepository,
	}
}

// GetMetricValue возвращает значение метрики по имени и типу
func (s *GetMetricValueService) GetMetricValue(req *handlers.GetMetricValueRequest) (string, error) {
	// Получаем текущее значение метрики
	currentValue, err := s.metricRepository.Get(req.Type, req.Name)
	if err != nil {
		// Возвращаем ошибку, если метрика не найдена
		return getValueEmptyString, &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: "metric not found",
		}
	}

	return currentValue, nil
}
