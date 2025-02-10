package services

import (
	"context"

	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"go-metrics-alerting/pkg/logger"
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

	// Получаем текущее значение метрики с использованием контекста
	currentValue, err := s.MetricRepository.Get(ctx, string(req.Type), req.Name)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"type": req.Type,
			"name": req.Name,
		}).Error("Failed to get metric value, setting default value to 0")
		currentValue = "0"
	}

	var value string

	// Обработка типа метрики
	switch req.Type {
	case types.Counter:
		logger.Logger.WithFields(logrus.Fields{
			"current_value": currentValue,
		}).Info("Processing Counter metric")

		newVal, _ := strconv.ParseInt(req.Value, 10, 64)
		curVal, _ := strconv.ParseInt(currentValue, 10, 64)
		newVal += curVal
		value = strconv.FormatInt(newVal, 10)

	case types.Gauge:
		logger.Logger.WithFields(logrus.Fields{
			"current_value": currentValue,
		}).Info("Processing Gauge metric")

		newVal, _ := strconv.ParseFloat(req.Value, 64)
		value = strconv.FormatFloat(newVal, 'f', -1, 64)
	}

	// Сохраняем новое значение с использованием контекста
	err = s.MetricRepository.Save(ctx, string(req.Type), req.Name, value)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"type":  req.Type,
			"name":  req.Name,
			"value": value,
		}).Error("Failed to save metric value")

		return &apierror.APIError{
			Code:    http.StatusInternalServerError,
			Message: errors.ErrValueNotSaved.Error(),
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
