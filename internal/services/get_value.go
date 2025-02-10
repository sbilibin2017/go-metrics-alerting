package services

import (
	"context"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
)

// Интерфейс для получения значения метрики
type MetricStorageGetter interface {
	// Метод для получения значения метрики по её типу и имени с использованием контекста
	Get(ctx context.Context, metricType string, metricName string) (string, error)
}

// Сервис для получения значения метрики
type GetMetricValueService struct {
	MetricRepository MetricStorageGetter // Репозиторий, который будет хранить метрики
}

// GetMetricValue возвращает значение метрики по имени и типу
func (s *GetMetricValueService) GetMetricValue(ctx context.Context, req *types.GetMetricValueRequest) (string, error) {
	// Получаем текущее значение метрики из репозитория с использованием контекста
	currentValue, err := s.MetricRepository.Get(ctx, req.Type, req.Name)
	if err != nil {
		// Возвращаем ошибку, если метрика не найдена
		return types.EmptyString, &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: errors.ErrMetricNotFound.Error(),
		}
	}

	// Возвращаем полученное значение метрики
	return currentValue, nil
}
