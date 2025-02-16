package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"

	"github.com/go-resty/resty/v2"
)

// Интерфейс HttpClient
type HttpClient interface {
	R() *resty.Request
}

// MetricAgentService - структура для сбора и отправки метрик.
type MetricAgentService struct {
	config    *configs.AgentConfig
	metricsCh chan types.MetricsRequest // Канал для передачи метрик
	client    HttpClient                // Интерфейс для отправки HTTP запросов
}

// NewMetricAgentService - создает новый сервис для сбора и отправки метрик.
func NewMetricAgentService(config *configs.AgentConfig, client HttpClient) *MetricAgentService {
	return &MetricAgentService{
		config:    config,
		metricsCh: make(chan types.MetricsRequest, 100),
		client:    client,
	}
}

// Start запускает процесс сбора и отправки метрик по расписанию.
func (s *MetricAgentService) Start() {
	tickerPoll := time.NewTicker(s.config.PollInterval * time.Second)
	tickerReport := time.NewTicker(s.config.ReportInterval * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	for {
		select {
		case <-tickerPoll.C:
			s.collectGauges()
			s.collectCounters()

		case <-tickerReport.C:
			s.sendMetrics()
			s.clearChannel()
		}
	}
}

// Метод для сбора метрик типа gauge
func (s *MetricAgentService) collectGauges() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	f := func(value float64) *float64 {
		return &value
	}

	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "Alloc", Value: f(float64(ms.Alloc))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "BuckHashSys", Value: f(float64(ms.BuckHashSys))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "Frees", Value: f(float64(ms.Frees))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "GCCPUFraction", Value: f(ms.GCCPUFraction)}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "HeapAlloc", Value: f(float64(ms.HeapAlloc))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "HeapIdle", Value: f(float64(ms.HeapIdle))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "HeapInuse", Value: f(float64(ms.HeapInuse))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "HeapObjects", Value: f(float64(ms.HeapObjects))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "HeapReleased", Value: f(float64(ms.HeapReleased))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "HeapSys", Value: f(float64(ms.HeapSys))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "NumGC", Value: f(float64(ms.NumGC))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "Sys", Value: f(float64(ms.Sys))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "TotalAlloc", Value: f(float64(ms.TotalAlloc))}
	s.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "RandomValue", Value: f(rand.Float64())}
}

// Метод для сбора счетчиков
func (s *MetricAgentService) collectCounters() {
	// Определение замыкания для инкрементации счётчика PollCount
	count := int64(0)
	counter := func() *int64 {
		count++
		return &count
	}

	// Создание метрики PollCount и отправка в канал
	s.metricsCh <- types.MetricsRequest{
		MType: "counter",
		ID:    "PollCount",
		Delta: counter(),
	}
}

// Метод для отправки метрик
func (s *MetricAgentService) sendMetrics() {
	for metric := range s.metricsCh {
		resp, err := s.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(metric).
			Post(fmt.Sprintf("%s/update/", s.config.Address))
		if err != nil {
			fmt.Println("Network error:", err)
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			fmt.Println("Unexpected status code:", resp.StatusCode())
		}
	}
}

// Метод для очистки канала
func (s *MetricAgentService) clearChannel() {
	for len(s.metricsCh) > 0 {
		<-s.metricsCh
	}
}
