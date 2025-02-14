package repositories

import (
	"errors"
	"strings"
)

const (
	keySeparator string = ":"
	emptyString  string = ""
)

var (
	ErrValueDoesNotExist error = errors.New("value does not exist")
)

// StorageEngine defines the interface for interacting with data storage.
type Storage interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Range(callback func(key, value string) bool)
}

// MetricRepository manages the saving and retrieval of metrics.
type MetricRepository struct {
	storage Storage
}

func NewMetricRepository(storage Storage) *MetricRepository {
	return &MetricRepository{storage: storage}
}

// Save stores a metric in the storage. Even if metricType or metricName is empty, it will be saved.
func (r *MetricRepository) Save(metricType, metricName, metricValue string) {
	key := strings.Join([]string{metricType, metricName}, keySeparator)
	r.storage.Set(key, metricValue)
}

// Get retrieves a metric by its type and name. It returns an error if the key is invalid.
func (r *MetricRepository) Get(metricType, metricName string) (string, error) {
	key := strings.Join([]string{metricType, metricName}, keySeparator)
	value, exists := r.storage.Get(key)
	if !exists {
		return emptyString, ErrValueDoesNotExist
	}
	return value, nil
}

// GetAll retrieves all valid metrics from the storage.
func (r *MetricRepository) GetAll() [][]string {
	var allMetrics [][]string
	r.storage.Range(func(key, value string) bool {
		parts := strings.Split(key, keySeparator)
		if len(parts) != 2 || parts[0] == emptyString || parts[1] == emptyString {
			return true
		}
		allMetrics = append(allMetrics, append(parts, value))
		return true
	})
	return allMetrics
}
