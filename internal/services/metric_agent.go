package services

import (
	"fmt"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/logger"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
)

// APIClient interface to be mocked
type APIClientEngine interface {
	R() *resty.Request
}

// MetricAgentService struct
type MetricAgentService struct {
	APIClient      APIClientEngine
	PollInterval   time.Duration
	ReportInterval time.Duration
	MetricChannel  chan types.UpdateMetricValueRequest
	Shutdown       chan os.Signal
	Address        string
}

// Start starts the metric agent service
func (s *MetricAgentService) Start() {
	signal.Notify(s.Shutdown, syscall.SIGINT, syscall.SIGTERM)
	tickerPoll := time.NewTicker(s.PollInterval)
	tickerReport := time.NewTicker(s.ReportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	var pollCount int64

	for {
		select {
		case <-tickerPoll.C:
			logger.Logger.Info("Collecting metrics...")

			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			pollCount++

			metrics := []types.UpdateMetricValueRequest{
				{Type: string(types.Gauge), Name: "Alloc", Value: fmt.Sprintf("%d", ms.Alloc)},
				{Type: string(types.Gauge), Name: "BuckHashSys", Value: fmt.Sprintf("%d", ms.BuckHashSys)},
				{Type: string(types.Gauge), Name: "Frees", Value: fmt.Sprintf("%d", ms.Frees)},
				{Type: string(types.Gauge), Name: "GCCPUFraction", Value: fmt.Sprintf("%f", ms.GCCPUFraction)},
				{Type: string(types.Gauge), Name: "HeapAlloc", Value: fmt.Sprintf("%d", ms.HeapAlloc)},
				{Type: string(types.Gauge), Name: "HeapIdle", Value: fmt.Sprintf("%d", ms.HeapIdle)},
				{Type: string(types.Gauge), Name: "HeapInuse", Value: fmt.Sprintf("%d", ms.HeapInuse)},
				{Type: string(types.Gauge), Name: "HeapObjects", Value: fmt.Sprintf("%d", ms.HeapObjects)},
				{Type: string(types.Gauge), Name: "HeapReleased", Value: fmt.Sprintf("%d", ms.HeapReleased)},
				{Type: string(types.Gauge), Name: "HeapSys", Value: fmt.Sprintf("%d", ms.HeapSys)},
				{Type: string(types.Gauge), Name: "NumGC", Value: fmt.Sprintf("%d", ms.NumGC)},
				{Type: string(types.Gauge), Name: "Sys", Value: fmt.Sprintf("%d", ms.Sys)},
				{Type: string(types.Gauge), Name: "TotalAlloc", Value: fmt.Sprintf("%d", ms.TotalAlloc)},
				{Type: string(types.Gauge), Name: "RandomValue", Value: fmt.Sprintf("%f", rand.Float64())},
				{Type: string(types.Counter), Name: "PollCount", Value: fmt.Sprintf("%d", pollCount)},
			}

			for _, metric := range metrics {
				logger.Logger.Debugf("Collected metric: %s %s = %s", metric.Type, metric.Name, metric.Value)
				s.MetricChannel <- metric
			}

		case <-tickerReport.C:
			logger.Logger.Info("Reporting metrics...")

			// Запуск отправки метрик в цикле
			for {
				select {
				case metric := <-s.MetricChannel:
					logger.Logger.Debugf("Sending metric: %s %s = %s", metric.Type, metric.Name, metric.Value)

					// Отправка по пути
					s.APIClient.R().Post(fmt.Sprintf("%s/update/%s/%s/%s", s.Address, metric.Type, metric.Name, metric.Value))

				default:
					// Прерываем выполнение текущего цикла, если нет метрик для отправки
					return
				}
			}

		case <-s.Shutdown:
			logger.Logger.Info("Shutting down gracefully...")
			return
		}
	}
}
