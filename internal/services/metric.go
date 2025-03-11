package services

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/types"
)

// SaveMetricRepository is an interface for saving a single metric
type SaveMetricRepository interface {
	SaveMetric(ctx context.Context, metric *types.Metrics) bool
}

// SaveMetricsRepository is an interface for saving multiple metrics
type SaveMetricsRepository interface {
	SaveMetrics(ctx context.Context, metrics []*types.Metrics) bool
}

// GetMetricByIdRepository is an interface for retrieving a metric by its ID
type GetMetricByIdRepository interface {
	GetMetricByID(ctx context.Context, id types.MetricID) *types.Metrics
}

// GetMetricsByIdsRepository is an interface for retrieving multiple metrics by their IDs
type GetMetricsByIdsRepository interface {
	GetMetricsByIDs(ctx context.Context, ids []types.MetricID) map[types.MetricID]*types.Metrics
}

// ListMetricsRepository is an interface for listing all metrics
type ListMetricsRepository interface {
	ListMetrics(ctx context.Context) []*types.Metrics
}

type MetricRepository interface {
	SaveMetricRepository
	SaveMetricsRepository
	GetMetricByIdRepository
	GetMetricsByIdsRepository
	ListMetricsRepository
}

var (
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
	ErrSaveMetric            = errors.New("failed to save metric")
	ErrMetricNotFound        = errors.New("metric not found")
)

type MetricService struct {
	*UpdateMetricService
	*UpdatesMetricService
	*GetMetricService
	*GetAllMetricsService
}

// NewMetricService creates a new instance of MetricService with initialized sub-services
func NewMetricService(repo MetricRepository) *MetricService {
	return &MetricService{
		UpdateMetricService:  &UpdateMetricService{repo: repo, getRepo: repo},
		UpdatesMetricService: &UpdatesMetricService{repo: repo, getRepo: repo},
		GetMetricService:     &GetMetricService{repo: repo},
		GetAllMetricsService: &GetAllMetricsService{repo: repo},
	}
}

// UpdateMetricService handles updating a single metric
type UpdateMetricService struct {
	repo    SaveMetricRepository
	getRepo GetMetricByIdRepository
}

// UpdateMetric updates a single metric based on its type (gauge or counter)
func (s *UpdateMetricService) UpdateMetric(ctx context.Context, metric *types.Metrics) (*types.Metrics, error) {
	existingMetric := s.getRepo.GetMetricByID(ctx, metric.MetricID)
	updater, ok := getUpdaterByType(metric.Type)
	if !ok {
		return nil, ErrUnsupportedMetricType
	}
	metric = updater(existingMetric, metric)
	ok = s.repo.SaveMetric(ctx, metric)
	if !ok {
		return nil, ErrSaveMetric
	}
	return metric, nil
}

func getUpdaterByType(
	tp types.MetricType,
) (func(existingMetric *types.Metrics, newMetric *types.Metrics) *types.Metrics, bool) {
	switch tp {
	case types.Counter:
		return updateCounter, true
	case types.Gauge:
		return updateGauge, true
	default:
		return nil, false
	}
}

func updateCounter(existingMetric *types.Metrics, newMetric *types.Metrics) *types.Metrics {
	if existingMetric != nil {
		*newMetric.Delta += *existingMetric.Delta
		return newMetric
	}
	return newMetric
}

func updateGauge(_ *types.Metrics, newMetric *types.Metrics) *types.Metrics {
	return newMetric
}

// UpdatesMetricService handles updating multiple metrics
type UpdatesMetricService struct {
	repo    SaveMetricsRepository
	getRepo GetMetricsByIdsRepository
}

// UpdateMetrics updates multiple metrics based on their types
func (s *UpdatesMetricService) UpdateMetrics(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error) {
	ids := make([]types.MetricID, len(metrics))
	for idx, m := range metrics {
		ids[idx] = m.MetricID
	}

	existingMetrics := s.getRepo.GetMetricsByIDs(ctx, ids)

	for idx, metric := range metrics {
		updater, ok := getUpdaterByType(metric.Type)
		if !ok {
			return nil, ErrUnsupportedMetricType
		}
		metrics[idx] = updater(existingMetrics[metric.MetricID], metric)
	}

	ok := s.repo.SaveMetrics(ctx, metrics)

	if !ok {
		return nil, ErrSaveMetric
	}

	return metrics, nil
}

// GetMetricService handles retrieving a single metric
type GetMetricService struct {
	repo GetMetricByIdRepository
}

// GetMetric retrieves a metric by its ID
func (s *GetMetricService) GetMetric(ctx context.Context, id types.MetricID) (*types.Metrics, error) {
	metric := s.repo.GetMetricByID(ctx, id)
	if metric == nil {
		return nil, ErrMetricNotFound
	}
	return metric, nil
}

// GetAllMetricsService handles retrieving all metrics
type GetAllMetricsService struct {
	repo ListMetricsRepository
}

// GetAllMetrics retrieves all stored metrics
func (s *GetAllMetricsService) GetAllMetrics(ctx context.Context) []*types.Metrics {
	metrics := s.repo.ListMetrics(ctx)
	return metrics
}
