package agent

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type MType string

const (
	Gauge   MType = "gauge"
	Counter MType = "counter"
)

type UpdateMetricsRequest struct {
	ID    string   `json:"id"`
	MType MType    `json:"mtype"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type AgentConfig struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`    // Интервал опроса, дефолтное значение 2 секунды
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"` // Интервал отчётов, дефолтное значение 10 секунд
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
}

func StartAgent(signalCh chan os.Signal, config *AgentConfig, client *resty.Client) {
	log.Printf("Starting agent with config: Address=%s, PollInterval=%v, ReportInterval=%v", config.Address, config.PollInterval, config.ReportInterval)

	metricsCh := make(chan UpdateMetricsRequest, 100)

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

func collectMetrics(metricsCh chan UpdateMetricsRequest) {
	metrics := append(collectGaugeMetrics(), collectCounterMetrics()()...)
	log.Printf("Collected %d metrics", len(metrics))
	for _, metric := range metrics {
		metricsCh <- metric
	}
}

func sendMetrics(metricsCh chan UpdateMetricsRequest, client *resty.Client, address string) {
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

func sendMetric(metric UpdateMetricsRequest, client *resty.Client, address string) {
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

func collectGaugeMetrics() []UpdateMetricsRequest {
	newGaugeMetric := func(id string, value interface{}) UpdateMetricsRequest {
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

		return UpdateMetricsRequest{
			MType: Gauge,
			ID:    id,
			Value: floatValue,
		}
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return []UpdateMetricsRequest{
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
func collectCounterMetrics() func() []UpdateMetricsRequest {
	var pollCount int64
	return func() []UpdateMetricsRequest {
		pollCount++
		delta := pollCount
		return []UpdateMetricsRequest{
			{
				MType: Counter,
				ID:    "PollCount",
				Delta: &delta,
			},
		}
	}
}
