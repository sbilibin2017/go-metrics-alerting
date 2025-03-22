package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"os"
	"sync"
)

type MetricFileRepository struct {
	file *os.File
	c    *configs.ServerConfig
	mu   sync.Mutex
}

// NewMetricFileRepository creates a new instance of MetricFileRepository with the provided file and ServerConfig.
// This now overwrites the file on every save.
func NewMetricFileRepository(c *configs.ServerConfig) (*MetricFileRepository, error) {
	// Open the file for writing and truncate it if it already exists
	file, err := os.OpenFile(c.FileStoragePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return &MetricFileRepository{
		file: file, // Directly use *os.File
		c:    c,
	}, nil
}

// SaveMetrics saves a list of metrics to a file, overwriting the existing content.
func (mr *MetricFileRepository) SaveMetrics(ctx context.Context, metrics []*types.Metrics) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	// Clear the file content by opening it with truncation mode
	// Re-open the file in truncation mode each time we write metrics to overwrite it
	file, err := os.OpenFile(mr.c.FileStoragePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	// Convert metrics to JSON and write to the file
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			return fmt.Errorf("failed to marshal metric: %v", err)
		}

		// Writing with a newline for each metric
		_, err = file.Write(append(data, '\n'))
		if err != nil {
			return fmt.Errorf("failed to write metric to file: %v", err)
		}
	}

	return nil
}

// FilterMetricsByTypeAndID filters metrics by their IDs and returns matching metrics.
func (mr *MetricFileRepository) FilterMetricsByTypeAndID(ctx context.Context, metricIDs []types.MetricID) ([]*types.Metrics, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	// Open the file for reading
	file, err := os.Open(mr.c.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer file.Close()

	var matchingMetrics []*types.Metrics // Slice to store matching metrics
	scanner := bufio.NewScanner(file)

	// Read the file line by line and unmarshal JSON into metrics
	for scanner.Scan() {
		var metric types.Metrics
		err := json.Unmarshal(scanner.Bytes(), &metric)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metric: %v", err)
		}

		// Check if the metric ID (ID + Type) matches any of the provided metric IDs
		for _, id := range metricIDs {
			// Comparing metricID (ID + Type) with the provided MetricID
			if metric.ID == id.ID && metric.Type == id.Type {
				// Add matching metric to the slice
				matchingMetrics = append(matchingMetrics, &metric)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Return the matching metrics
	return matchingMetrics, nil
}

// ListMetrics lists all the metrics stored in the file.
func (mr *MetricFileRepository) ListMetrics(ctx context.Context) ([]*types.Metrics, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	// Open the file for reading
	file, err := os.Open(mr.c.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer file.Close()

	var metrics []*types.Metrics // Slice to store all metrics
	scanner := bufio.NewScanner(file)

	// Read the file line by line and unmarshal JSON into metrics
	for scanner.Scan() {
		var metric types.Metrics
		err := json.Unmarshal(scanner.Bytes(), &metric)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metric: %v", err)
		}

		// Add the metric to the slice
		metrics = append(metrics, &metric)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Return all the metrics
	return metrics, nil
}
