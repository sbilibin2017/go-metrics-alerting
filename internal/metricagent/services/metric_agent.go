package services

import (
	"fmt"
	"go-metrics-alerting/pkg/logger"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const chBuffSize int = 100

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

type UpdateMetricValueRequest struct {
	Type  MetricType
	Name  string
	Value string
}

// MetricFacade интерфейс для фасада работы с метриками
// UpdateMetric обновляет значение метрики
type MetricFacade interface {
	UpdateMetric(metricType string, metricName string, metricValue string) error
}

// MetricAgentService сервис для сбора и отправки метрик
type MetricAgentService struct {
	metricFacade   MetricFacade
	metricChannel  chan UpdateMetricValueRequest
	pollInterval   time.Duration
	reportInterval time.Duration
	shutdown       chan os.Signal
}

// NewMetricAgentService создает новый экземпляр MetricAgentService
func NewMetricAgentService(
	metricFacade MetricFacade,
	pollInterval time.Duration,
	reportInterval time.Duration,
) *MetricAgentService {
	return &MetricAgentService{
		metricFacade:   metricFacade,
		metricChannel:  make(chan UpdateMetricValueRequest, chBuffSize),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		shutdown:       make(chan os.Signal, 1), // Канал для сигналов
	}
}

// Start запускает процесс сбора и отправки метрик
func (s *MetricAgentService) Start() {
	// Перехват сигналов для graceful shutdown
	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM)

	tickerPoll := time.NewTicker(s.pollInterval)
	tickerReport := time.NewTicker(s.reportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	collectors := map[string]func() []UpdateMetricValueRequest{
		"gauge":   collectGaugeMetrics,
		"counter": collectCounterMetrics,
	}

	for {
		select {
		case <-tickerPoll.C:
			// Логирование начала сбора метрик
			logger.Logger.Info("Collecting metrics...")

			for _, collector := range collectors {
				metrics := collector()
				for _, metric := range metrics {
					// Логирование каждой метрики
					logger.Logger.Debugf("Collected metric: %s %s = %s", metric.Type, metric.Name, metric.Value)
					s.metricChannel <- metric
				}
			}

		case <-tickerReport.C:
			// Логирование начала отправки метрик
			logger.Logger.Info("Reporting metrics...")

			// Используем флаг, чтобы корректно выйти из внутреннего цикла
			for {
				select {
				case metric := <-s.metricChannel:
					// Логирование отправляемой метрики
					logger.Logger.Debugf("Sending metric: %s %s = %s", metric.Type, metric.Name, metric.Value)

					if err := s.metricFacade.UpdateMetric(string(metric.Type), metric.Name, metric.Value); err != nil {
						// Логирование ошибки при отправке
						logger.Logger.Errorf("Error updating %s metric %s: %v", metric.Type, metric.Name, err)
					}
				default:
					// Завершаем внутренний цикл
					return
				}
			}

		case <-s.shutdown:
			// Сигнал завершения работы
			logger.Logger.Info("Shutting down gracefully...") // Логирование завершения работы
			return
		}
	}
}

// collectGaugeMetrics собирает метрики типа Gauge
func collectGaugeMetrics() []UpdateMetricValueRequest {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	metrics := []UpdateMetricValueRequest{
		{Type: GaugeType, Name: "Alloc", Value: fmt.Sprintf("%d", ms.Alloc)},
		{Type: GaugeType, Name: "BuckHashSys", Value: fmt.Sprintf("%d", ms.BuckHashSys)},
		{Type: GaugeType, Name: "Frees", Value: fmt.Sprintf("%d", ms.Frees)},
		{Type: GaugeType, Name: "GCCPUFraction", Value: fmt.Sprintf("%f", ms.GCCPUFraction)},
		{Type: GaugeType, Name: "HeapAlloc", Value: fmt.Sprintf("%d", ms.HeapAlloc)},
		{Type: GaugeType, Name: "HeapIdle", Value: fmt.Sprintf("%d", ms.HeapIdle)},
		{Type: GaugeType, Name: "HeapInuse", Value: fmt.Sprintf("%d", ms.HeapInuse)},
		{Type: GaugeType, Name: "HeapObjects", Value: fmt.Sprintf("%d", ms.HeapObjects)},
		{Type: GaugeType, Name: "HeapReleased", Value: fmt.Sprintf("%d", ms.HeapReleased)},
		{Type: GaugeType, Name: "HeapSys", Value: fmt.Sprintf("%d", ms.HeapSys)},
		{Type: GaugeType, Name: "NumGC", Value: fmt.Sprintf("%d", ms.NumGC)},
		{Type: GaugeType, Name: "Sys", Value: fmt.Sprintf("%d", ms.Sys)},
		{Type: GaugeType, Name: "TotalAlloc", Value: fmt.Sprintf("%d", ms.TotalAlloc)},
		{Type: GaugeType, Name: "RandomValue", Value: fmt.Sprintf("%f", rand.Float64())},
	}

	// Логирование всех собранных метрик типа Gauge
	for _, metric := range metrics {
		logger.Logger.Debugf("Collected Gauge Metric: %s = %s", metric.Name, metric.Value)
	}

	return metrics
}

var pollCount int64

// collectCounterMetrics собирает метрики типа Counter
func collectCounterMetrics() []UpdateMetricValueRequest {
	pollCount++
	metrics := []UpdateMetricValueRequest{
		{Type: CounterType, Name: "PollCount", Value: fmt.Sprintf("%d", pollCount)},
	}

	// Логирование собранной метрики типа Counter
	for _, metric := range metrics {
		logger.Logger.Debugf("Collected Counter Metric: %s = %s", metric.Name, metric.Value)
	}

	return metrics
}
