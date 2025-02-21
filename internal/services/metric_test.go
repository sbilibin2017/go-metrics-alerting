package services_test

import (
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Моки для Saver и Getter
type mockSaver struct {
	store map[string]*types.Metrics
}

func (m *mockSaver) Save(key string, value *types.Metrics) bool {
	m.store[key] = value
	return true
}

type mockGetter struct {
	store map[string]*types.Metrics
}

func (m *mockGetter) Get(key string) (*types.Metrics, bool) {
	metric, exists := m.store[key]
	return metric, exists
}

func TestUpdateGaugeMetricService(t *testing.T) {
	saver := &mockSaver{store: make(map[string]*types.Metrics)}
	getter := &mockGetter{store: saver.store}
	service := services.NewUpdateGaugeMetricService(saver, getter)

	value := 42.5
	metric := &types.Metrics{ID: "gauge1", MType: types.Gauge, Value: &value}
	updatedMetric, ok := service.Update(metric)

	// Use testify assertions
	assert.True(t, ok, "Expected successful update")
	assert.NotNil(t, updatedMetric.Value, "Expected value to be set")
	assert.Equal(t, value, *updatedMetric.Value, "Expected value to match")
}

func TestUpdateCounterMetricService(t *testing.T) {
	saver := &mockSaver{store: make(map[string]*types.Metrics)}
	getter := &mockGetter{store: saver.store}
	service := services.NewUpdateCounterMetricService(saver, getter)

	value := int64(10)
	metric := &types.Metrics{ID: "counter1", MType: types.Counter, Delta: &value}
	updatedMetric, ok := service.Update(metric)

	// Use testify assertions
	assert.True(t, ok, "Expected successful update")
	assert.NotNil(t, updatedMetric.Delta, "Expected delta to be set")
	assert.Equal(t, value, *updatedMetric.Delta, "Expected delta to match")

	// Проверка накопления значения
	newValue := int64(5)
	metric2 := &types.Metrics{ID: "counter1", MType: types.Counter, Delta: &newValue}
	updatedMetric, _ = service.Update(metric2)

	expectedSum := int64(15)
	assert.Equal(t, expectedSum, *updatedMetric.Delta, "Expected accumulated delta to match")
}

func TestUpdateMetricService(t *testing.T) {
	saver := &mockSaver{store: make(map[string]*types.Metrics)}
	getter := &mockGetter{store: saver.store}
	gaugeService := services.NewUpdateGaugeMetricService(saver, getter)
	counterService := services.NewUpdateCounterMetricService(saver, getter)
	service := services.NewUpdateMetricService(gaugeService, counterService)

	// Тест gauge
	gaugeValue := 55.5
	gaugeMetric := &types.Metrics{ID: "gauge2", MType: types.Gauge, Value: &gaugeValue}
	updatedGauge, ok := service.Update(gaugeMetric)

	// Use testify assertions
	assert.True(t, ok, "Expected successful update")
	assert.NotNil(t, updatedGauge.Value, "Expected value to be set for gauge")
	assert.Equal(t, gaugeValue, *updatedGauge.Value, "Expected gauge value to match")

	// Тест counter
	counterValue := int64(20)
	counterMetric := &types.Metrics{ID: "counter2", MType: types.Counter, Delta: &counterValue}
	updatedCounter, ok := service.Update(counterMetric)

	// Use testify assertions
	assert.True(t, ok, "Expected successful update")
	assert.NotNil(t, updatedCounter.Delta, "Expected delta to be set for counter")
	assert.Equal(t, counterValue, *updatedCounter.Delta, "Expected counter delta to match")
}

func TestUpdateExistingGaugeMetric(t *testing.T) {
	saver := &mockSaver{store: make(map[string]*types.Metrics)}
	getter := &mockGetter{store: saver.store}
	service := services.NewUpdateGaugeMetricService(saver, getter)

	// Setup an existing metric
	existingValue := 42.5
	existingMetric := &types.Metrics{ID: "gauge1", MType: types.Gauge, Value: &existingValue}
	saver.Save(existingMetric.ID, existingMetric) // Save the initial metric

	// New metric with the same ID, but different value
	newValue := 55.5
	metric := &types.Metrics{ID: "gauge1", MType: types.Gauge, Value: &newValue}

	// Perform the update
	updatedMetric, ok := service.Update(metric)

	// Use testify assertions
	assert.True(t, ok, "Expected successful update")
	assert.NotNil(t, updatedMetric.Value, "Expected value to be set")
	assert.Equal(t, newValue, *updatedMetric.Value, "Expected updated value to match")

	// Verify that the existing metric was updated correctly
	storedMetric, exists := saver.store[existingMetric.ID]
	assert.True(t, exists, "Expected the metric to be saved in the store")
	assert.Equal(t, newValue, *storedMetric.Value, "Expected stored metric value to match updated value")
}

func TestUpdateInvalidMetricType(t *testing.T) {
	saver := &mockSaver{store: make(map[string]*types.Metrics)}
	getter := &mockGetter{store: saver.store}
	service := services.NewUpdateMetricService(
		services.NewUpdateGaugeMetricService(saver, getter),
		services.NewUpdateCounterMetricService(saver, getter),
	)

	// Creating a metric with an unsupported type (assuming there is no handler for this type)
	invalidMetric := &types.Metrics{
		ID:    "invalid1",
		MType: "unsupported_type", // Assuming this type is not handled by your services
	}

	// Perform the update (this should trigger the default case and return nil, false)
	updatedMetric, ok := service.Update(invalidMetric)

	// Use testify assertions
	assert.False(t, ok, "Expected update to fail for unsupported metric type")
	assert.Nil(t, updatedMetric, "Expected returned metric to be nil for unsupported metric type")
}
