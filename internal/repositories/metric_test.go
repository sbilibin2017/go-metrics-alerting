package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockStorage implements the Storage interface for testing purposes.
type MockStorage struct {
	data map[string]string
}

func (m *MockStorage) Set(key, value string) {
	m.data[key] = value
}

func (m *MockStorage) Get(key string) (string, bool) {
	value, exists := m.data[key]
	return value, exists
}

func (m *MockStorage) Range(callback func(key, value string) bool) {
	for key, value := range m.data {
		if !callback(key, value) {
			break
		}
	}
}

func TestNewMetricRepository(t *testing.T) {
	// Создаем мок хранилища
	mockStorage := &MockStorage{data: make(map[string]string)}

	// Создаем репозиторий с помощью конструктора
	repo := NewMetricRepository(mockStorage)

	// Проверка, что репозиторий был инициализирован и использует переданное хранилище
	assert.NotNil(t, repo)
	assert.Equal(t, mockStorage, repo.storage)
}

// TestSave tests saving a metric to the storage.
func TestSave(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	repo := &MetricRepository{storage: mockStorage}

	// Save a valid metric
	repo.Save("cpu", "usage", "75")

	// Assert that the metric is saved correctly
	value, exists := mockStorage.Get("cpu:usage")
	assert.True(t, exists)
	assert.Equal(t, "75", value)
}

// TestSaveWithEmptyTypeOrName tests saving a metric with empty metricType or metricName.
func TestSaveWithEmptyTypeOrName(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	repo := &MetricRepository{storage: mockStorage}

	// Save a metric with an empty type
	repo.Save("", "usage", "50")
	// Save a metric with an empty name
	repo.Save("cpu", "", "30")

	// Assert that the metrics with empty type or name are still saved
	// The keys will be ":usage" and "cpu:"
	value, exists := mockStorage.Get(":usage")
	assert.True(t, exists)
	assert.Equal(t, "50", value)

	value, exists = mockStorage.Get("cpu:")
	assert.True(t, exists)
	assert.Equal(t, "30", value)
}

// TestGet tests retrieving a metric from the storage.
func TestGet(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	repo := &MetricRepository{storage: mockStorage}

	// Save a metric
	repo.Save("cpu", "usage", "75")

	// Test getting the saved metric
	value, err := repo.Get("cpu", "usage")
	assert.NoError(t, err)
	assert.Equal(t, "75", value)

	// Test getting a non-existing metric
	_, err = repo.Get("cpu", "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, ErrValueDoesNotExist, err)
}

// TestGetAll tests retrieving all metrics from the storage.
func TestGetAll(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	repo := &MetricRepository{storage: mockStorage}

	// Save a few metrics
	repo.Save("cpu", "usage", "75")
	repo.Save("disk", "usage", "50")
	repo.Save("", "empty", "100")   // Saving with empty metricType
	repo.Save("network", "", "200") // Saving with empty metricName

	// Get all metrics
	allMetrics := repo.GetAll()

	// Assert that we have 2 valid metrics (ignoring the invalid ones)
	assert.Len(t, allMetrics, 2)

	// Check that the valid metrics are returned correctly
	assert.Contains(t, allMetrics, []string{"cpu", "usage", "75"})
	assert.Contains(t, allMetrics, []string{"disk", "usage", "50"})
}

// TestGetAllWithInvalidKeys tests GetAll with invalid key formats.
func TestGetAllWithInvalidKeys(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	repo := &MetricRepository{storage: mockStorage}

	// Save some invalid metrics
	repo.Save("", "empty", "100")           // Invalid key with empty type
	repo.Save("network", "", "200")         // Invalid key with empty name
	repo.Save("disk:extra", "usage", "300") // Invalid format with more than one separator

	// Get all metrics
	allMetrics := repo.GetAll()

	// Assert that no invalid metrics are included
	assert.Len(t, allMetrics, 0)
}
