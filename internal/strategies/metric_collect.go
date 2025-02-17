package strategies

import (
	"go-metrics-alerting/internal/types"
	"math/rand"
	"runtime"
)

// GaugeCollector - структура для сбора метрик типа gauge
type GaugeCollector struct{}

// Collect - метод для сбора метрик типа gauge
func (c *GaugeCollector) Collect() []types.MetricsRequest {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	f := func(value float64) *float64 {
		return &value
	}

	return []types.MetricsRequest{
		{MType: "gauge", ID: "Alloc", Value: f(float64(ms.Alloc))},
		{MType: "gauge", ID: "BuckHashSys", Value: f(float64(ms.BuckHashSys))},
		{MType: "gauge", ID: "Frees", Value: f(float64(ms.Frees))},
		{MType: "gauge", ID: "GCCPUFraction", Value: f(ms.GCCPUFraction)},
		{MType: "gauge", ID: "HeapAlloc", Value: f(float64(ms.HeapAlloc))},
		{MType: "gauge", ID: "HeapIdle", Value: f(float64(ms.HeapIdle))},
		{MType: "gauge", ID: "HeapInuse", Value: f(float64(ms.HeapInuse))},
		{MType: "gauge", ID: "HeapObjects", Value: f(float64(ms.HeapObjects))},
		{MType: "gauge", ID: "HeapReleased", Value: f(float64(ms.HeapReleased))},
		{MType: "gauge", ID: "HeapSys", Value: f(float64(ms.HeapSys))},
		{MType: "gauge", ID: "NumGC", Value: f(float64(ms.NumGC))},
		{MType: "gauge", ID: "Sys", Value: f(float64(ms.Sys))},
		{MType: "gauge", ID: "TotalAlloc", Value: f(float64(ms.TotalAlloc))},
		{MType: "gauge", ID: "RandomValue", Value: f(rand.Float64())},
	}
}

// CounterCollector - структура для сбора метрик типа counter
type CounterCollector struct {
	count int64
}

// Collect - метод для сбора метрик типа counter
func (c *CounterCollector) Collect() []types.MetricsRequest {
	c.count++
	return []types.MetricsRequest{
		{MType: "counter", ID: "PollCount", Delta: &c.count},
	}
}
