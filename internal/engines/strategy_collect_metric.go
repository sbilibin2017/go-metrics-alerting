package engines

import (
	"go-metrics-alerting/internal/types"
	"math/rand"
	"runtime"
)

// MemoryStatsProvider интерфейс для получения статистики о памяти
type MemoryStatsProvider interface {
	ReadMemStats(ms *runtime.MemStats)
}

// Структура для получения реальной статистики
type RealMemoryStatsProvider struct{}

// Реализация метода ReadMemStats для структуры RealMemoryStatsProvider
func (r *RealMemoryStatsProvider) ReadMemStats(ms *runtime.MemStats) {
	runtime.ReadMemStats(ms)
}

// MetricCollectionStrategyEngineInterface интерфейс для стратегий сбора метрик
type MetricCollectionStrategyEngineInterface interface {
	CollectMetrics() []interface{}
}

// GaugeCollectionStrategyEngine собирает метрики типа Gauge
type GaugeCollectionStrategyEngine struct {
	memStatsProvider MemoryStatsProvider
}

// Конструктор для GaugeCollectionStrategyEngine
func NewGaugeCollectionStrategyEngine(memStatsProvider MemoryStatsProvider) *GaugeCollectionStrategyEngine {
	return &GaugeCollectionStrategyEngine{memStatsProvider: memStatsProvider}
}

// CollectMetrics собирает метрики о памяти
func (g *GaugeCollectionStrategyEngine) CollectMetrics() []interface{} {
	// Проверка на nil перед использованием интерфейса
	if g.memStatsProvider == nil {
		// Возвращаем пустой срез, если интерфейс не инициализирован
		return nil
	}

	var memStats runtime.MemStats
	g.memStatsProvider.ReadMemStats(&memStats)

	metrics := []interface{}{
		createMetric(string(types.GaugeType), "Alloc", FormatFloat(float64(memStats.Alloc))),
		createMetric(string(types.GaugeType), "BuckHashSys", FormatFloat(float64(memStats.BuckHashSys))),
		createMetric(string(types.GaugeType), "Frees", FormatFloat(float64(memStats.Frees))),
		createMetric(string(types.GaugeType), "GCCPUFraction", FormatFloat(float64(memStats.GCCPUFraction))),
		createMetric(string(types.GaugeType), "GCSys", FormatFloat(float64(memStats.GCSys))),
		createMetric(string(types.GaugeType), "HeapAlloc", FormatFloat(float64(memStats.HeapAlloc))),
		createMetric(string(types.GaugeType), "HeapIdle", FormatFloat(float64(memStats.HeapIdle))),
		createMetric(string(types.GaugeType), "HeapInuse", FormatFloat(float64(memStats.HeapInuse))),
		createMetric(string(types.GaugeType), "HeapObjects", FormatFloat(float64(memStats.HeapObjects))),
		createMetric(string(types.GaugeType), "HeapReleased", FormatFloat(float64(memStats.HeapReleased))),
		createMetric(string(types.GaugeType), "HeapSys", FormatFloat(float64(memStats.HeapSys))),
		createMetric(string(types.GaugeType), "LastGC", FormatFloat(float64(memStats.LastGC))),
		createMetric(string(types.GaugeType), "Lookups", FormatFloat(float64(memStats.Lookups))),
		createMetric(string(types.GaugeType), "MCacheInuse", FormatFloat(float64(memStats.MCacheInuse))),
		createMetric(string(types.GaugeType), "MCacheSys", FormatFloat(float64(memStats.MCacheSys))),
		createMetric(string(types.GaugeType), "MSpanInuse", FormatFloat(float64(memStats.MSpanInuse))),
		createMetric(string(types.GaugeType), "MSpanSys", FormatFloat(float64(memStats.MSpanSys))),
		createMetric(string(types.GaugeType), "Mallocs", FormatFloat(float64(memStats.Mallocs))),
		createMetric(string(types.GaugeType), "NextGC", FormatFloat(float64(memStats.NextGC))),
		createMetric(string(types.GaugeType), "NumForcedGC", FormatFloat(float64(memStats.NumForcedGC))),
		createMetric(string(types.GaugeType), "NumGC", FormatFloat(float64(memStats.NumGC))),
		createMetric(string(types.GaugeType), "OtherSys", FormatFloat(float64(memStats.OtherSys))),
		createMetric(string(types.GaugeType), "PauseTotalNs", FormatFloat(float64(memStats.PauseTotalNs))),
		createMetric(string(types.GaugeType), "StackInuse", FormatFloat(float64(memStats.StackInuse))),
		createMetric(string(types.GaugeType), "StackSys", FormatFloat(float64(memStats.StackSys))),
		createMetric(string(types.GaugeType), "Sys", FormatFloat(float64(memStats.Sys))),
		createMetric(string(types.GaugeType), "TotalAlloc", FormatFloat(float64(memStats.TotalAlloc))),
		createMetric(string(types.GaugeType), "RandomValue", rand.Float64()*100.0),
	}

	return metrics
}

// CounterCollectionStrategyEngine собирает метрики типа Counter
type CounterCollectionStrategyEngine struct {
	pollCount int64
}

// CollectMetrics для Counter метрик
func (c *CounterCollectionStrategyEngine) CollectMetrics() []interface{} {
	metrics := []interface{}{
		createMetric(string(types.CounterType), "PollCount", FormatInt(c.pollCount)),
	}
	c.pollCount++
	return metrics
}

// helper function для форматирования метрики
func createMetric(metricType, name string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Type":  metricType,
		"Name":  name,
		"Value": value,
	}
}
