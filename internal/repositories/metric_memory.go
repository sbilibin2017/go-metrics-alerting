package repositories

import (
	"context"
	"go-metrics-alerting/internal/types"
)

type MetricMemoryRepository struct {
	data map[types.MetricID]*types.Metrics
}

// NewMetricMemoryRepository creates a new instance of MetricMemoryRepository.
func NewMetricMemoryRepository() *MetricMemoryRepository {
	return &MetricMemoryRepository{
		data: make(map[types.MetricID]*types.Metrics), // Initialize the map
	}
}

// SaveMetrics saves a list of metrics in the in-memory storage.
func (mr *MetricMemoryRepository) SaveMetrics(ctx context.Context, metrics []*types.Metrics) error {
	for _, metric := range metrics {
		// Create a MetricID for the key
		metricID := types.MetricID{ID: metric.ID, Type: metric.Type}

		// Store the metric in memory using MetricID as the key
		mr.data[metricID] = metric
	}
	return nil
}

// FilterMetricsByTypeAndId filters metrics by their IDs and types, and returns matching metrics.
func (mr *MetricMemoryRepository) FilterMetricsByTypeAndId(ctx context.Context, metricIDs []types.MetricID) ([]*types.Metrics, error) {
	var result []*types.Metrics

	// Iterate through the list of MetricID objects
	for _, metricID := range metricIDs {
		// Retrieve the metric from memory using MetricID as the key
		metric, exists := mr.data[metricID]
		if exists {
			result = append(result, metric)
		}
	}

	return result, nil
}

// ListMetrics lists all metrics stored in the in-memory storage.
func (mr *MetricMemoryRepository) ListMetrics(ctx context.Context) ([]*types.Metrics, error) {
	var result []*types.Metrics

	// Iterate through the in-memory map and collect all metrics
	for _, metric := range mr.data {
		result = append(result, metric)
	}

	return result, nil
}
