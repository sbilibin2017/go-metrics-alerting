package services

import (
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/errors" // Import the errors package
	"go-metrics-alerting/internal/types"
	"net/http"
)

// MetricServiceInterface defines the methods that the MetricService should implement
type MetricServiceInterface interface {
	// UpdateMetric updates a metric based on the provided request.
	UpdateMetric(req types.UpdateMetricRequest) errors.ApiErrorInterface

	// GetMetric retrieves the value of a metric based on the provided request.
	GetMetric(req types.GetMetricRequest) (string, errors.ApiErrorInterface)

	// GetAllMetrics retrieves all metrics.
	GetAllMetrics() []types.UpdateMetricRequest
}

// MetricService manages metrics with strategies, storage engine, and key engine
type MetricService struct {
	storageEngine   engines.StorageEngineInterface
	strategyEngines map[types.MetricType]engines.StrategyUpdateEngineInterface
	keyEngine       engines.KeyEngineInterface
}

// NewMetricService creates a new MetricService with a storage engine, key engine, and strategies map
func NewMetricService(
	storageEngine engines.StorageEngineInterface,
	strategyEngines map[types.MetricType]engines.StrategyUpdateEngineInterface,
	keyEngine engines.KeyEngineInterface,
) *MetricService {
	return &MetricService{
		storageEngine:   storageEngine,
		strategyEngines: strategyEngines,
		keyEngine:       keyEngine,
	}
}

// UpdateMetric updates a metric using the appropriate strategy and stores the value
func (m *MetricService) UpdateMetric(req types.UpdateMetricRequest) errors.ApiErrorInterface {
	if err := m.isValidMetricType(req.Type); err != nil {
		return err
	}
	if err := m.isValidMetricExists(req.Type); err != nil {
		return err
	}
	if err := m.isValidMetricName(req.Name); err != nil {
		return err
	}
	if err := m.isValidMetricValue(req.Value); err != nil {
		return err
	}
	encodedKey := m.keyEngine.Encode(req.Type, req.Name)
	currentValue, exists := m.storageEngine.Get(encodedKey)
	if !exists {
		currentValue = "0"
	}
	strategy, ok := m.strategyEngines[types.MetricType(req.Type)]
	if !ok {
		return &errors.ApiError{StatusCode: http.StatusBadRequest, Message: "strategy not found"}
	}
	updatedValue, err := strategy.Update(currentValue, req.Value)
	if err != nil {
		return &errors.ApiError{StatusCode: http.StatusBadRequest, Message: "invalid metric value"}
	}
	m.storageEngine.Set(encodedKey, updatedValue)
	return nil
}

// GetMetric retrieves a metric value from storage
func (m *MetricService) GetMetric(req types.GetMetricRequest) (string, errors.ApiErrorInterface) {
	if err := m.isValidMetricType(req.Type); err != nil {
		return "", err
	}
	if err := m.isValidMetricName(req.Name); err != nil {
		return "", err
	}
	encodedKey := m.keyEngine.Encode(req.Type, req.Name)
	value, exists := m.storageEngine.Get(encodedKey)
	if !exists {
		return "", &errors.ApiError{StatusCode: http.StatusNotFound, Message: "metric not found"}
	}
	return value, nil
}

// GetAllMetrics returns all metrics from storage
func (m *MetricService) GetAllMetrics() []types.UpdateMetricRequest {
	var metrics []types.UpdateMetricRequest
	for pair := range m.storageEngine.Generate() {
		metricType, name, err := m.keyEngine.Decode(pair[0])
		if err != nil {
			continue
		}
		metrics = append(metrics, types.UpdateMetricRequest{
			Name:  name,
			Value: pair[1],
			Type:  metricType,
		})
	}
	return metrics
}

// validateMetricType checks if the metric type is provided and valid
func (m *MetricService) isValidMetricType(metricType string) errors.ApiErrorInterface {
	if metricType == "" {
		return &errors.ApiError{StatusCode: http.StatusBadRequest, Message: "metric type is required"}
	}
	return nil
}

// validateMetricName checks if the metric name is provided
func (m *MetricService) isValidMetricName(metricName string) errors.ApiErrorInterface {
	if metricName == "" {
		return &errors.ApiError{StatusCode: http.StatusBadRequest, Message: "metric name is required"}
	}
	return nil
}

// validateMetricValue checks if the metric value is provided
func (m *MetricService) isValidMetricValue(metricValue string) errors.ApiErrorInterface {
	if metricValue == "" {
		return &errors.ApiError{StatusCode: http.StatusBadRequest, Message: "metric value is required"}
	}
	return nil
}

// validateMetricExists checks if the metric type exists in the strategy engines
func (m *MetricService) isValidMetricExists(metricType string) errors.ApiErrorInterface {
	if _, exists := m.strategyEngines[types.MetricType(metricType)]; !exists {
		return &errors.ApiError{StatusCode: http.StatusBadRequest, Message: "invalid metric type"}
	}
	return nil
}

// Ensure ChannelEngine implements the ChannelEngineInterface at compile time.
var _ MetricServiceInterface = &MetricService{}
