package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/types"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Start запускает процесс сбора и отправки метрик по расписанию.
func StartMetricAgent(config *configs.AgentConfig, client *resty.Client) {
	metricsCh := make(chan types.UpdateMetricsRequest, 100)

	tickerPoll := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	// Логирование старта
	logger.Logger.Info("Metric collection started",
		zap.Int("poll_interval", int(config.PollInterval)),
		zap.Int("report_interval", int(config.ReportInterval)))

	// Бесконечный цикл для сбора и отправки метрик
	for {
		select {
		case <-tickerPoll.C:
			// Сбор метрик
			collectMetrics(metricsCh)
		case <-tickerReport.C:
			// Отправка метрик
			sendMetrics(metricsCh, client, config.Address)
		}
	}
}

// collectMetrics - функция для сбора метрик и отправки их в канал.
func collectMetrics(metricsCh chan types.UpdateMetricsRequest) {
	gaugeMetrics := collectGaugeMetrics()
	counterMetrics := collectCounterMetrics()

	// Отправляем метрики в канал
	for _, metric := range append(gaugeMetrics, counterMetrics...) {
		metricsCh <- metric
	}

	logger.Logger.Debug("Collected metrics",
		zap.Int("gauge_metrics_count", len(gaugeMetrics)),
		zap.Int("counter_metrics_count", len(counterMetrics)))
}

// sendMetrics - функция для отправки метрик через REST API.
func sendMetrics(metricsCh chan types.UpdateMetricsRequest, client *resty.Client, address string) {
	for metric := range metricsCh {
		err := sendMetric(metric, client, address)
		if err != nil {
			logger.Logger.Error("Error sending metric",
				zap.String("metric_id", metric.ID),
				zap.Error(err))
		} else {
			logger.Logger.Debug("Metric sent successfully",
				zap.String("metric_id", metric.ID))
		}
	}
}

// sendMetric - функция для отправки одной метрики через REST API.
func sendMetric(metric types.UpdateMetricsRequest, client *resty.Client, address string) error {
	url := address + "/update/"

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metric).
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		return fmt.Errorf("server returned error: %d", resp.StatusCode())
	}

	logger.Logger.Debug("Metric sent",
		zap.String("metric_id", metric.ID),
		zap.Int("status_code", resp.StatusCode()))

	return nil
}

// collectGaugeMetrics - функция для сбора метрик типа gauge.
func collectGaugeMetrics() []types.UpdateMetricsRequest {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	f := func(value float64) *float64 {
		return &value
	}

	metrics := []types.UpdateMetricsRequest{
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

	logger.Logger.Debug("Collected gauge metrics", zap.Int("count", len(metrics)))
	return metrics
}

// collectCounterMetrics - функция для сбора метрик типа counter.
func collectCounterMetrics() []types.UpdateMetricsRequest {
	var count int64

	updateCount := func() {
		count++
	}

	updateCount()

	metrics := []types.UpdateMetricsRequest{
		{MType: "counter", ID: "PollCount", Delta: &count},
	}

	logger.Logger.Debug("Collected counter metrics", zap.Int64("PollCount", count))
	return metrics
}
