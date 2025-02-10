package repositories

import (
	"context"
	"errors"
	"go-metrics-alerting/pkg/logger" // Import the logger package

	"github.com/sirupsen/logrus"
)

// Константы для работы с метриками и ошибками
const (
	EmptyString = "" // Пустая строка
)

// Ошибки
var (
	ErrMetricNotFound  = errors.New("metric not found")     // Ошибка при отсутствии метрики
	ErrKeyDecodeFailed = errors.New("failed to decode key") // Ошибка при неудаче декодирования ключа
)

// StorageEngine defines the interface for interacting with data storage.
type StorageEngine interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Generate(ctx context.Context) <-chan []string
}

// KeyEngine defines the interface for encoding and decoding keys.
type KeyEngine interface {
	Encode(mt, mn string) string
	Decode(key string) (string, string, error)
}

// MetricRepository manages the saving and retrieval of metrics.
type MetricRepository struct {
	StorageEngine StorageEngine
	KeyEngine     KeyEngine
}

// Save stores a metric in the storage.
func (r *MetricRepository) Save(ctx context.Context, metricType, metricName, metricValue string) error {
	key := r.KeyEngine.Encode(metricType, metricName)

	// Log the start of saving a metric
	logger.Logger.WithFields(logrus.Fields{
		"metricType":  metricType,
		"metricName":  metricName,
		"metricValue": metricValue,
	}).Info("Saving metric")

	// Save the metric in storage
	err := r.StorageEngine.Set(ctx, key, metricValue)
	if err != nil {
		// Log the error in case of failure
		logger.Logger.WithFields(logrus.Fields{
			"key":   key,
			"error": err.Error(),
		}).Error("Failed to save metric")
		return err
	}

	// Log success
	logger.Logger.WithFields(logrus.Fields{
		"key": key,
	}).Info("Metric saved successfully")

	return nil
}

// Get retrieves a metric by its type and name.
func (r *MetricRepository) Get(ctx context.Context, metricType, metricName string) (string, error) {
	key := r.KeyEngine.Encode(metricType, metricName)

	// Log the start of retrieving a metric
	logger.Logger.WithFields(logrus.Fields{
		"metricType": metricType,
		"metricName": metricName,
	}).Info("Retrieving metric")

	// Get the metric from storage
	value, err := r.StorageEngine.Get(ctx, key)
	if err != nil {
		// Log the error in case of failure
		logger.Logger.WithFields(logrus.Fields{
			"key":   key,
			"error": err.Error(),
		}).Error("Failed to retrieve metric")
		return EmptyString, err
	}

	// Log success
	logger.Logger.WithFields(logrus.Fields{
		"key":   key,
		"value": value,
	}).Info("Metric retrieved successfully")

	return value, nil
}

// GetAll retrieves all metrics from the storage.
func (r *MetricRepository) GetAll(ctx context.Context) [][]string {
	allMetrics := [][]string{}

	// Log the start of retrieving all metrics
	logger.Logger.Info("Retrieving all metrics")

	// Using context for Generate method
	for kv := range r.StorageEngine.Generate(ctx) {
		mt, mn, err := r.KeyEngine.Decode(kv[0])
		if err != nil {
			// Log error in decoding key
			logger.Logger.WithFields(logrus.Fields{
				"key":   kv[0],
				"error": err.Error(),
			}).Warn("Failed to decode key")
			continue
		}

		// Log decoded metric
		logger.Logger.WithFields(logrus.Fields{
			"metricType":  mt,
			"metricName":  mn,
			"metricValue": kv[1],
		}).Info("Decoded metric successfully")

		allMetrics = append(allMetrics, []string{mt, mn, kv[1]})
	}

	// Log the end of the retrieval process
	logger.Logger.Info("All metrics retrieved successfully")
	return allMetrics
}
