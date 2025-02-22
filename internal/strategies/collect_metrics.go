package strategies

import (
	"go-metrics-alerting/internal/types"
	"math/rand"
	"runtime"
)

// GaugeMetricsCollector собирает Gauge метрики
type GaugeMetricsCollector struct{}

// NewGaugeMetricsCollector создает новый экземпляр GaugeMetricsCollector
func NewGaugeMetricsCollector() *GaugeMetricsCollector {
	return &GaugeMetricsCollector{}
}

// Collect собирает метрики типа Gauge и возвращает их в формате UpdateMetricBodyRequest
func (g *GaugeMetricsCollector) Collect() []*types.UpdateMetricBodyRequest {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return []*types.UpdateMetricBodyRequest{
		{MType: "gauge", ID: "Alloc", Value: float64Ptr(float64(ms.Alloc))},
		{MType: "gauge", ID: "BuckHashSys", Value: float64Ptr(float64(ms.BuckHashSys))},
		{MType: "gauge", ID: "Frees", Value: float64Ptr(float64(ms.Frees))},
		{MType: "gauge", ID: "GCCPUFraction", Value: float64Ptr(ms.GCCPUFraction)},
		{MType: "gauge", ID: "HeapAlloc", Value: float64Ptr(float64(ms.HeapAlloc))},
		{MType: "gauge", ID: "HeapIdle", Value: float64Ptr(float64(ms.HeapIdle))},
		{MType: "gauge", ID: "HeapInuse", Value: float64Ptr(float64(ms.HeapInuse))},
		{MType: "gauge", ID: "HeapObjects", Value: float64Ptr(float64(ms.HeapObjects))},
		{MType: "gauge", ID: "HeapReleased", Value: float64Ptr(float64(ms.HeapReleased))},
		{MType: "gauge", ID: "HeapSys", Value: float64Ptr(float64(ms.HeapSys))},
		{MType: "gauge", ID: "NumGC", Value: float64Ptr(float64(ms.NumGC))},
		{MType: "gauge", ID: "Sys", Value: float64Ptr(float64(ms.Sys))},
		{MType: "gauge", ID: "TotalAlloc", Value: float64Ptr(float64(ms.TotalAlloc))},
		{MType: "gauge", ID: "RandomValue", Value: float64Ptr(rand.Float64())},
	}
}

// CounterMetricsCollector собирает Counter метрики
type CounterMetricsCollector struct {
	pollCount int64
}

// NewCounterMetricsCollector создает новый экземпляр CounterMetricsCollector
func NewCounterMetricsCollector() *CounterMetricsCollector {
	return &CounterMetricsCollector{}
}

// Collect собирает метрики типа Counter и возвращает их в формате UpdateMetricBodyRequest
func (c *CounterMetricsCollector) Collect() []*types.UpdateMetricBodyRequest {
	c.pollCount++
	delta := c.pollCount
	return []*types.UpdateMetricBodyRequest{
		{MType: "counter", ID: "PollCount", Delta: int64Ptr(delta)},
	}
}

// Вспомогательная функция для создания указателя на float64
func float64Ptr(v float64) *float64 {
	return &v
}

// Вспомогательная функция для создания указателя на int64
func int64Ptr(v int64) *int64 {
	return &v
}
