package services

import (
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

// MetricCollector интерфейс для всех коллекционеров метрик
type MetricCollector interface {
	Collect() []types.UpdateMetricValueRequest
}

// MetricAgentService отвечает за сбор и отправку метрик.
type MetricAgentService struct {
	config     *configs.AgentConfig
	apiClient  *resty.Client
	metrics    []types.UpdateMetricValueRequest
	collectors map[types.MetricType]MetricCollector
}

// NewMetricAgentService создает и инициализирует новый объект MetricAgentService.
func NewMetricAgentService(config *configs.AgentConfig) *MetricAgentService {
	return &MetricAgentService{
		config:    config,
		apiClient: resty.New(),
		metrics:   make([]types.UpdateMetricValueRequest, 0),
		collectors: map[types.MetricType]MetricCollector{
			types.Gauge:   &GaugeCollector{},
			types.Counter: &CounterCollector{},
		},
	}
}

// Start запускает сервис сбора метрик.
func (s *MetricAgentService) Start() {
	tickerPoll := time.NewTicker(s.config.PollInterval)
	tickerReport := time.NewTicker(s.config.ReportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	for {
		select {
		case <-tickerPoll.C:
			s.collectMetrics()
		case <-tickerReport.C:
			s.sendMetrics()
		}
	}
}

// collectMetrics собирает метрики от всех коллекционеров.
func (s *MetricAgentService) collectMetrics() {
	for _, collector := range s.collectors {
		collected := collector.Collect()
		if collected != nil {
			s.metrics = append(s.metrics, collected...)
		}
	}
}

// sendMetrics отправляет переданные метрики через API и очищает срез метрик.
func (s *MetricAgentService) sendMetrics() {
	for _, metric := range s.metrics {
		url := fmt.Sprintf("%s/update/%s/%s/%s", s.config.Address, metric.Type, metric.Name, metric.Value)
		s.apiClient.R().Post(url)
	}
	s.metrics = []types.UpdateMetricValueRequest{}
}

// GaugeCollector собирает метрики типа Gauge.
type GaugeCollector struct{}

func (g *GaugeCollector) Collect() []types.UpdateMetricValueRequest {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	formatValue := func(value float64) string {
		return fmt.Sprintf("%f", value)
	}

	return []types.UpdateMetricValueRequest{
		{Type: types.Gauge, Name: "Alloc", Value: formatValue(float64(ms.Alloc))},
		{Type: types.Gauge, Name: "BuckHashSys", Value: formatValue(float64(ms.BuckHashSys))},
		{Type: types.Gauge, Name: "Frees", Value: formatValue(float64(ms.Frees))},
		{Type: types.Gauge, Name: "GCCPUFraction", Value: formatValue(ms.GCCPUFraction)},
		{Type: types.Gauge, Name: "HeapAlloc", Value: formatValue(float64(ms.HeapAlloc))},
		{Type: types.Gauge, Name: "HeapIdle", Value: formatValue(float64(ms.HeapIdle))},
		{Type: types.Gauge, Name: "HeapInuse", Value: formatValue(float64(ms.HeapInuse))},
		{Type: types.Gauge, Name: "HeapObjects", Value: formatValue(float64(ms.HeapObjects))},
		{Type: types.Gauge, Name: "HeapReleased", Value: formatValue(float64(ms.HeapReleased))},
		{Type: types.Gauge, Name: "HeapSys", Value: formatValue(float64(ms.HeapSys))},
		{Type: types.Gauge, Name: "NumGC", Value: formatValue(float64(ms.NumGC))},
		{Type: types.Gauge, Name: "Sys", Value: formatValue(float64(ms.Sys))},
		{Type: types.Gauge, Name: "TotalAlloc", Value: formatValue(float64(ms.TotalAlloc))},
		{Type: types.Gauge, Name: "RandomValue", Value: formatValue(rand.Float64())},
	}
}

// CounterCollector отслеживает счетчик событий.
type CounterCollector struct {
	count int64
}

func (c *CounterCollector) Collect() []types.UpdateMetricValueRequest {
	c.count++

	formatValue := func(value int64) string {
		return fmt.Sprintf("%d", value)
	}

	return []types.UpdateMetricValueRequest{
		{Type: types.Counter, Name: "PollCount", Value: formatValue(int64(c.count))},
	}
}
