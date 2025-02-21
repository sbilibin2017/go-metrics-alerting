package strategies

import (
	"go-metrics-alerting/internal/domain"
	"math/rand"
	"runtime"
)

// GaugeMetricsCollector собирает Gauge метрики
type GaugeMetricsCollector struct{}

// Collect собирает метрики типа Gauge и возвращает их в формате доменной модели Metric
func (g *GaugeMetricsCollector) Collect() []*domain.Metrics {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return []*domain.Metrics{
		{MType: domain.Gauge, ID: "Alloc", Value: float64Ptr(float64(ms.Alloc))},
		{MType: domain.Gauge, ID: "BuckHashSys", Value: float64Ptr(float64(ms.BuckHashSys))},
		{MType: domain.Gauge, ID: "Frees", Value: float64Ptr(float64(ms.Frees))},
		{MType: domain.Gauge, ID: "GCCPUFraction", Value: float64Ptr(ms.GCCPUFraction)},
		{MType: domain.Gauge, ID: "HeapAlloc", Value: float64Ptr(float64(ms.HeapAlloc))},
		{MType: domain.Gauge, ID: "HeapIdle", Value: float64Ptr(float64(ms.HeapIdle))},
		{MType: domain.Gauge, ID: "HeapInuse", Value: float64Ptr(float64(ms.HeapInuse))},
		{MType: domain.Gauge, ID: "HeapObjects", Value: float64Ptr(float64(ms.HeapObjects))},
		{MType: domain.Gauge, ID: "HeapReleased", Value: float64Ptr(float64(ms.HeapReleased))},
		{MType: domain.Gauge, ID: "HeapSys", Value: float64Ptr(float64(ms.HeapSys))},
		{MType: domain.Gauge, ID: "NumGC", Value: float64Ptr(float64(ms.NumGC))},
		{MType: domain.Gauge, ID: "Sys", Value: float64Ptr(float64(ms.Sys))},
		{MType: domain.Gauge, ID: "TotalAlloc", Value: float64Ptr(float64(ms.TotalAlloc))},
		{MType: domain.Gauge, ID: "RandomValue", Value: float64Ptr(rand.Float64())},
	}
}

// CounterMetricsCollector собирает Counter метрики
type CounterMetricsCollector struct {
	pollCount int64
}

// Collect собирает метрики типа Counter и возвращает их в формате доменной модели Metric
func (c *CounterMetricsCollector) Collect() []*domain.Metrics {
	c.pollCount++
	delta := c.pollCount
	return []*domain.Metrics{
		{MType: domain.Counter, ID: "PollCount", Delta: int64Ptr(delta)},
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
