package services

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/types"
)

type MetricRepository interface {
	SaveMetrics(ctx context.Context, metrics []*types.Metrics) error
	FilterMetricsByTypeAndID(ctx context.Context, metricIDs []types.MetricID) ([]*types.Metrics, error)
	ListMetrics(ctx context.Context) ([]*types.Metrics, error)
}

type MetricService struct {
	repo MetricRepository
}

func NewMetricService(repo MetricRepository) *MetricService {
	return &MetricService{repo: repo}
}

// UpdatesMetric updates the metrics and returns the updated metrics.
func (s *MetricService) UpdatesMetric(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error) {
	var metricIDs []types.MetricID
	for _, metric := range metrics {
		metricIDs = append(metricIDs, types.MetricID{ID: metric.ID, Type: metric.Type})
	}

	// Filter existing metrics
	existingMetrics, err := s.repo.FilterMetricsByTypeAndID(ctx, metricIDs)
	if err != nil {
		return nil, err
	}

	// Create a map for existing metrics by MetricID
	metricMap := make(map[types.MetricID]*types.Metrics)
	for _, metric := range existingMetrics {
		metricMap[types.MetricID{ID: metric.ID, Type: metric.Type}] = metric
	}

	// Update existing metrics or add new ones
	for _, metric := range metrics {
		existingMetric, exists := metricMap[types.MetricID{ID: metric.ID, Type: metric.Type}]
		if exists {
			// Update the existing metric based on type
			switch metric.Type {
			case string(types.Gauge):
				existingMetric.Value = metric.Value
			case string(types.Counter):
				*existingMetric.Delta += *metric.Delta
			}
		} else {
			// Add the new metric to the map
			metricMap[types.MetricID{ID: metric.ID, Type: metric.Type}] = metric
		}
	}

	// Prepare the list of updated metrics
	var updatedMetrics []*types.Metrics
	for _, metric := range metricMap {
		updatedMetrics = append(updatedMetrics, metric)
	}

	// Save updated metrics to the repository
	if err := s.repo.SaveMetrics(ctx, updatedMetrics); err != nil {
		return nil, err
	}

	return updatedMetrics, nil
}

// GetMetricByTypeAndID fetches a metric by its ID and type.
func (s *MetricService) GetMetricByTypeAndID(ctx context.Context, id types.MetricID) (*types.Metrics, error) {
	// Filter metrics by ID and type
	metrics, err := s.repo.FilterMetricsByTypeAndID(ctx, []types.MetricID{id})
	if err != nil {
		return nil, err
	}

	// Check if the metric is found
	if len(metrics) == 0 {
		return nil, ErrMetricNotFound
	}

	return metrics[0], nil
}

// ListAllMetrics returns all the metrics in the repository.
func (s *MetricService) ListAllMetrics(ctx context.Context) ([]*types.Metrics, error) {
	// List all metrics
	metrics, err := s.repo.ListMetrics(ctx)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// Helper for logging errors related to not finding a metric.
var ErrMetricNotFound = errors.New("not found")
