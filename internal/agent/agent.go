package agent

import (
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartAgent(signalCh chan os.Signal, config *configs.AgentConfig, client *resty.Client) {
	log.Printf("Starting agent with config: Address=%s, PollInterval=%v, ReportInterval=%v", config.Address, config.PollInterval, config.ReportInterval)

	metricsCh := make(chan types.UpdateMetricsRequest, 100)

	tickerPoll := time.NewTicker(config.PollInterval)
	tickerReport := time.NewTicker(config.ReportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	for {
		select {
		case <-signalCh:
			log.Println("Received shutdown signal. Stopping agent...")
			return
		case <-tickerPoll.C:
			log.Println("Collecting metrics...")
			collectMetrics(metricsCh)
		case <-tickerReport.C:
			log.Println("Sending metrics...")
			sendMetrics(metricsCh, client, config.Address)
		}
	}
}

func collectMetrics(metricsCh chan types.UpdateMetricsRequest) {
	metrics := append(collectGaugeMetrics(), collectCounterMetrics()()...)
	log.Printf("Collected %d metrics", len(metrics))
	for _, metric := range metrics {
		metricsCh <- metric
	}
}

func sendMetrics(metricsCh chan types.UpdateMetricsRequest, client *resty.Client, address string) {
	for {
		select {
		case metric := <-metricsCh:
			log.Printf("Sending metric: %v", metric)
			sendMetric(metric, client, address)
		default:
			return
		}
	}
}

func sendMetric(metric types.UpdateMetricsRequest, client *resty.Client, address string) {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = fmt.Sprintf("http://%s", address)
	}
	log.Printf("Sending metric to address: %s", address)
	_, err := client.R().SetHeader("Content-Type", "application/json").SetBody(metric).Post(address + "/update/")
	if err != nil {
		log.Printf("Error sending metric: %v", err)
	} else {
		log.Println("Metric sent successfully")
	}
}

func collectGaugeMetrics() []types.UpdateMetricsRequest {
	newGaugeMetric := func(id string, value interface{}) types.UpdateMetricsRequest {
		var floatValue *float64
		switch v := value.(type) {
		case uint64:
			f := float64(v)
			floatValue = &f
		case uint32:
			f := float64(v)
			floatValue = &f
		case float64:
			floatValue = &v
		}

		return types.UpdateMetricsRequest{
			MType: types.Gauge,
			ID:    id,
			Value: floatValue,
		}
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return []types.UpdateMetricsRequest{
		newGaugeMetric("Alloc", ms.Alloc),
		newGaugeMetric("BuckHashSys", ms.BuckHashSys),
		newGaugeMetric("Frees", ms.Frees),
		newGaugeMetric("GCCPUFraction", ms.GCCPUFraction),
		newGaugeMetric("HeapAlloc", ms.HeapAlloc),
		newGaugeMetric("HeapIdle", ms.HeapIdle),
		newGaugeMetric("HeapInuse", ms.HeapInuse),
		newGaugeMetric("HeapObjects", ms.HeapObjects),
		newGaugeMetric("HeapReleased", ms.HeapReleased),
		newGaugeMetric("HeapSys", ms.HeapSys),
		newGaugeMetric("NumGC", ms.NumGC),
		newGaugeMetric("Sys", ms.Sys),
		newGaugeMetric("TotalAlloc", ms.TotalAlloc),
		newGaugeMetric("RandomValue", rand.Float64()),
	}
}

// collectCounterMetrics обновляет переменную pollCount и возвращает замыкание для дальнейшего обновления
func collectCounterMetrics() func() []types.UpdateMetricsRequest {
	var pollCount int64
	return func() []types.UpdateMetricsRequest {
		pollCount++
		delta := pollCount
		return []types.UpdateMetricsRequest{
			{
				MType: types.Counter,
				ID:    "PollCount",
				Delta: &delta,
			},
		}
	}
}
