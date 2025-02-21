package strategies

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGaugeMetricsCollector проверяет сбор метрик типа Gauge
func TestGaugeMetricsCollector(t *testing.T) {
	collector := &GaugeMetricsCollector{}
	metrics := collector.Collect()

	assert.NotEmpty(t, metrics, "Gauge metrics should not be empty")

	expectedMetrics := map[string]bool{
		"Alloc": true, "BuckHashSys": true, "Frees": true,
		"GCCPUFraction": true, "HeapAlloc": true, "HeapIdle": true,
		"HeapInuse": true, "HeapObjects": true, "HeapReleased": true,
		"HeapSys": true, "NumGC": true, "Sys": true,
		"TotalAlloc": true, "RandomValue": true,
	}

	for _, metric := range metrics {
		assert.NotNil(t, metric.Value, "Gauge metric value should not be nil")
		assert.Nil(t, metric.Delta, "Gauge metric delta should be nil")
		assert.Equal(t, domain.Gauge, metric.MType, "Metric type should be Gauge")
		assert.Contains(t, expectedMetrics, metric.ID, "Unexpected metric ID")
	}
}

// TestCounterMetricsCollector проверяет сбор метрик типа Counter
func TestCounterMetricsCollector(t *testing.T) {
	collector := &CounterMetricsCollector{}

	firstMetrics := collector.Collect()
	assert.Len(t, firstMetrics, 1, "Counter metrics should have exactly one metric")
	assert.Equal(t, "PollCount", firstMetrics[0].ID, "Counter metric ID should be PollCount")
	assert.NotNil(t, firstMetrics[0].Delta, "Counter metric delta should not be nil")
	assert.Nil(t, firstMetrics[0].Value, "Counter metric value should be nil")
	assert.Equal(t, domain.Counter, firstMetrics[0].MType, "Metric type should be Counter")

	secondMetrics := collector.Collect()
	assert.Greater(t, *secondMetrics[0].Delta, *firstMetrics[0].Delta, "Counter metric delta should increase")
}
