package services

import (
	"context"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"go-metrics-alerting/pkg/logger" // Импортируем логгер
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// MetricRepository — интерфейс хранилища метрик.
type MetricRepository interface {
	// Метод для получения метрики
	Get(ctx context.Context, metricType, metricName string) (string, error)
	// Метод для сохранения метрики
	Save(ctx context.Context, metricType, metricName, value string) error
}

// UpdateMetricValueService — сервис для обновления значений метрик.
type UpdateMetricValueService struct {
	MetricRepository MetricRepository
}

// UpdateMetricValue обновляет значение метрики.
func (s *UpdateMetricValueService) UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	// Логируем начало работы функции
	logger.Logger.WithFields(logrus.Fields{
		"type":  req.Type,
		"name":  req.Name,
		"value": req.Value,
	}).Info("Started updating metric value")

	// Проверка наличия типа метрики
	if req.Type == types.EmptyString {
		logger.Logger.Warn("Metric type is required")
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: errors.ErrInvalidMetricType.Error(),
		}
	}

	// Проверка имени метрики
	if req.Name == types.EmptyString {
		logger.Logger.Warn("Invalid metric name")
		return &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: errors.ErrInvalidMetricName.Error(),
		}
	}

	// Проверка значения метрики
	if req.Value == types.EmptyString {
		logger.Logger.Warn("Invalid metric value")
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: errors.ErrInvalidMetricValue.Error(),
		}
	}

	// Получаем текущее значение метрики с использованием контекста
	currentValue, err := s.MetricRepository.Get(ctx, req.Type, req.Name)
	if err != nil {
		// Если ошибка, установим значение по умолчанию "0"
		logger.Logger.WithFields(logrus.Fields{
			"type": req.Type,
			"name": req.Name,
		}).Error("Failed to get metric value, setting default value to 0")
		currentValue = "0"
	}

	var value string

	// Обработка типа метрики
	switch req.Type {
	case string(types.Counter):
		// Логируем обработку типа Counter
		logger.Logger.WithFields(logrus.Fields{
			"current_value": currentValue,
		}).Info("Processing Counter metric")

		newVal, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			logger.Logger.Error("Failed to convert new value to int64")
			return &apierror.APIError{
				Code:    http.StatusBadRequest,
				Message: errors.ErrInvalidValue.Error(),
			}
		}

		// Преобразуем текущие и новые значения в int64.
		curVal, err := strconv.ParseInt(currentValue, 10, 64)
		if err != nil {
			logger.Logger.Error("Failed to convert current value to int64")
			return &apierror.APIError{
				Code:    http.StatusBadRequest,
				Message: errors.ErrInvalidValue.Error(),
			}
		}

		// Суммируем и сохраняем новое значение
		newVal += curVal

		// Присваиваем значение для сохранения
		value = strconv.FormatInt(newVal, 10)

	case string(types.Gauge):
		// Логируем обработку типа Gauge
		logger.Logger.WithFields(logrus.Fields{
			"current_value": currentValue,
		}).Info("Processing Gauge metric")

		// Преобразуем новое значение в float64
		newVal, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			logger.Logger.Error("Failed to convert new value to float64")
			return &apierror.APIError{
				Code:    http.StatusBadRequest,
				Message: errors.ErrInvalidValue.Error(),
			}
		}

		// Присваиваем значение для сохранения
		value = strconv.FormatFloat(newVal, 'f', -1, 64)

	default:
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: errors.ErrUnsupportedMetricType.Error(),
		}
	}

	// Сохраняем новое значение с использованием контекста
	err = s.MetricRepository.Save(ctx, req.Type, req.Name, value)
	if err != nil {
		// Логируем ошибку с соответствующими полями
		logger.Logger.WithFields(logrus.Fields{
			"type":  req.Type,
			"name":  req.Name,
			"value": value,
		}).Error("Failed to save metric value")

		// Возвращаем ошибку через APIError
		return &apierror.APIError{
			Code:    http.StatusInternalServerError,
			Message: "metric value is not saved", // Custom message for Save failure
		}
	}

	// Успешно сохранили значение метрики
	logger.Logger.WithFields(logrus.Fields{
		"type":  req.Type,
		"name":  req.Name,
		"value": value,
	}).Info("Successfully updated metric value")

	return nil
}
