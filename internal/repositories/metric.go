package repositories

import (
	"context"
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

const (
	MetricEmptyString string = ""
)

// MetricRepository manages the saving and retrieval of metrics.
type MetricRepository struct {
	StorageEngine StorageEngine
	KeyEngine     KeyEngine
}

// Save stores a metric in the storage.
func (r *MetricRepository) Save(ctx context.Context, metricType, metricName, metricValue string) error {
	key := r.KeyEngine.Encode(metricType, metricName)
	err := r.StorageEngine.Set(ctx, key, metricValue)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a metric by its type and name.
func (r *MetricRepository) Get(ctx context.Context, metricType, metricName string) (string, error) {
	key := r.KeyEngine.Encode(metricType, metricName)
	value, err := r.StorageEngine.Get(ctx, key)
	if err != nil {
		return MetricEmptyString, err
	}
	return value, nil
}

// GetAll retrieves all metrics from the storage.
func (r *MetricRepository) GetAll(ctx context.Context) [][]string {
	allMetrics := [][]string{}
	for kv := range r.StorageEngine.Generate(ctx) {
		mt, mn, err := r.KeyEngine.Decode(kv[0])
		if err != nil {
			continue
		}
		allMetrics = append(allMetrics, []string{mt, mn, kv[1]})
	}
	return allMetrics
}
