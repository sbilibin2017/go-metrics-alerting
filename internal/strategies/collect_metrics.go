package strategies

import (
	"fmt"
	"go-metrics-alerting/internal/domain"
	"math/rand"
	"runtime"
)

// GaugeMetricsCollector собирает Gauge метрики
type GaugeMetricsCollector struct{}

// Collect собирает метрики типа Gauge и возвращает их в формате доменной модели Metric
func (g *GaugeMetricsCollector) Collect() []*domain.Metric {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return []*domain.Metric{
		{
			MType: domain.Gauge,
			ID:    "Alloc",
			Value: fmt.Sprintf("%f", float64(ms.Alloc)),
		},
		{
			MType: domain.Gauge,
			ID:    "BuckHashSys",
			Value: fmt.Sprintf("%f", float64(ms.BuckHashSys)),
		},
		{
			MType: domain.Gauge,
			ID:    "Frees",
			Value: fmt.Sprintf("%f", float64(ms.Frees)),
		},
		{
			MType: domain.Gauge,
			ID:    "GCCPUFraction",
			Value: fmt.Sprintf("%f", ms.GCCPUFraction),
		},
		{
			MType: domain.Gauge,
			ID:    "HeapAlloc",
			Value: fmt.Sprintf("%f", float64(ms.HeapAlloc)),
		},
		{
			MType: domain.Gauge,
			ID:    "HeapIdle",
			Value: fmt.Sprintf("%f", float64(ms.HeapIdle)),
		},
		{
			MType: domain.Gauge,
			ID:    "HeapInuse",
			Value: fmt.Sprintf("%f", float64(ms.HeapInuse)),
		},
		{
			MType: domain.Gauge,
			ID:    "HeapObjects",
			Value: fmt.Sprintf("%f", float64(ms.HeapObjects)),
		},
		{
			MType: domain.Gauge,
			ID:    "HeapReleased",
			Value: fmt.Sprintf("%f", float64(ms.HeapReleased)),
		},
		{
			MType: domain.Gauge,
			ID:    "HeapSys",
			Value: fmt.Sprintf("%f", float64(ms.HeapSys)),
		},
		{
			MType: domain.Gauge,
			ID:    "NumGC",
			Value: fmt.Sprintf("%f", float64(ms.NumGC)),
		},
		{
			MType: domain.Gauge,
			ID:    "Sys",
			Value: fmt.Sprintf("%f", float64(ms.Sys)),
		},
		{
			MType: domain.Gauge,
			ID:    "TotalAlloc",
			Value: fmt.Sprintf("%f", float64(ms.TotalAlloc)),
		},
		{
			MType: domain.Gauge,
			ID:    "RandomValue",
			Value: fmt.Sprintf("%f", rand.Float64()),
		},
	}
}

// CounterMetricsCollector собирает Counter метрики
type CounterMetricsCollector struct {
	pollCount int64
}

// Collect собирает метрики типа Counter и возвращает их в формате доменной модели Metric
func (c *CounterMetricsCollector) Collect() []*domain.Metric {
	c.pollCount++
	delta := c.pollCount
	return []*domain.Metric{
		{
			MType: domain.Counter,
			ID:    "PollCount",
			Value: fmt.Sprintf("%d", delta),
		},
	}
}
