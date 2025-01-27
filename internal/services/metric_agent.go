package services

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/facades"
	"go-metrics-alerting/internal/types"
)

// MetricAgentService собирает метрики и отправляет их на сервер
type MetricAgentService struct {
	config       configs.AgentConfigInterface
	metricFacade facades.MetricFacadeInterface
	metricsChan  chan *types.UpdateMetricRequest // Канал для передачи метрик
	pollCount    int64                           // Счётчик обновлений метрик
}

// NewMetricAgentService создает новый агент
func NewMetricAgentService(
	config configs.AgentConfigInterface,
	metricFacade facades.MetricFacadeInterface,
) *MetricAgentService {
	return &MetricAgentService{
		config:       config,
		metricFacade: metricFacade,
		metricsChan:  make(chan *types.UpdateMetricRequest, 100), // Канал с буфером 100
	}
}

// Start запускает агент
func (a *MetricAgentService) Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Горутин для сбора метрик
	go func() {
		defer wg.Done()
		a.runCollectMetrics()
	}()

	// Горутин для отправки метрик
	go func() {
		defer wg.Done()
		a.runSendMetrics()
	}()

	wg.Wait()
}

// runCollectMetrics собирает метрики с интервалом PollInterval
func (a *MetricAgentService) runCollectMetrics() {
	for {
		a.collectMetrics()
		time.Sleep(a.config.GetPollInterval()) // pollInterval
	}
}

// collectMetrics собирает метрики из runtime и отправляет их в канал
func (a *MetricAgentService) collectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Создаем карту метрик
	metrics := map[string]interface{}{
		"Alloc":         float64(memStats.Alloc),        // Тип gauge (float64)
		"BuckHashSys":   float64(memStats.BuckHashSys),  // Тип gauge (float64)
		"Frees":         float64(memStats.Frees),        // Тип gauge (float64)
		"GCCPUFraction": memStats.GCCPUFraction,         // Тип gauge (float64)
		"GCSys":         float64(memStats.GCSys),        // Тип gauge (float64)
		"HeapAlloc":     float64(memStats.HeapAlloc),    // Тип gauge (float64)
		"HeapIdle":      float64(memStats.HeapIdle),     // Тип gauge (float64)
		"HeapInuse":     float64(memStats.HeapInuse),    // Тип gauge (float64)
		"HeapObjects":   float64(memStats.HeapObjects),  // Тип gauge (float64)
		"HeapReleased":  float64(memStats.HeapReleased), // Тип gauge (float64)
		"HeapSys":       float64(memStats.HeapSys),      // Тип gauge (float64)
		"LastGC":        float64(memStats.LastGC),       // Тип gauge (float64)
		"Lookups":       float64(memStats.Lookups),      // Тип gauge (float64)
		"MCacheInuse":   float64(memStats.MCacheInuse),  // Тип gauge (float64)
		"MCacheSys":     float64(memStats.MCacheSys),    // Тип gauge (float64)
		"MSpanInuse":    float64(memStats.MSpanInuse),   // Тип gauge (float64)
		"MSpanSys":      float64(memStats.MSpanSys),     // Тип gauge (float64)
		"Mallocs":       float64(memStats.Mallocs),      // Тип gauge (float64)
		"NextGC":        float64(memStats.NextGC),       // Тип gauge (float64)
		"NumForcedGC":   float64(memStats.NumForcedGC),  // Тип gauge (float64)
		"NumGC":         float64(memStats.NumGC),        // Тип gauge (float64)
		"OtherSys":      float64(memStats.OtherSys),     // Тип gauge (float64)
		"PauseTotalNs":  float64(memStats.PauseTotalNs), // Тип gauge (float64)
		"StackInuse":    float64(memStats.StackInuse),   // Тип gauge (float64)
		"StackSys":      float64(memStats.StackSys),     // Тип gauge (float64)
		"Sys":           float64(memStats.Sys),          // Тип gauge (float64)
		"TotalAlloc":    float64(memStats.TotalAlloc),   // Тип gauge (float64)

		// Дополнительные метрики
		"PollCount":   a.pollCount,          // Тип counter (int64)
		"RandomValue": rand.Float64() * 100, // Случайное число от 0 до 100
	}

	// Увеличиваем PollCount после каждого обновления
	a.pollCount++

	// Отправляем метрики в канал
	for name, value := range metrics {
		var metricType string
		var valueStr string

		// Определяем тип метрики и форматируем значение
		if name == "PollCount" {
			metricType = "counter"              // PollCount будет типа counter
			valueStr = fmt.Sprintf("%d", value) // Форматируем как int64
		} else {
			metricType = "gauge" // Все остальные метрики типа gauge

			// Проверяем тип значения и форматируем его соответственно
			switch v := value.(type) {
			case int64:
				// Для int64 метрик (например, PollCount)
				valueStr = fmt.Sprintf("%d", v)
			case float64:
				// Для float64 метрик (например, из runtime)
				valueStr = fmt.Sprintf("%f", v)
			default:
				log.Printf("Unknown value type: %T\n", v)
				continue
			}
		}

		// Создаем структуру UpdateMetricRequest
		updateRequest := &types.UpdateMetricRequest{
			Type:  metricType,
			Name:  name,
			Value: valueStr,
		}

		a.metricsChan <- updateRequest
	}

	log.Printf("Collected %d metrics, PollCount: %d\n", len(metrics), a.pollCount)
}

// runSendMetrics отправляет метрики с интервалом ReportInterval
func (a *MetricAgentService) runSendMetrics() {
	for {
		a.sendMetrics()
		time.Sleep(a.config.GetReportInterval()) // reportInterval
	}
}

// sendMetrics отправляет метрики из канала на сервер
func (a *MetricAgentService) sendMetrics() {
	// Собираем все метрики из канала
	var metricsToSend []*types.UpdateMetricRequest

	// Читаем все метрики из канала
	for {
		select {
		case metricRequest := <-a.metricsChan:
			metricsToSend = append(metricsToSend, metricRequest)
		default:
			// Если канал пуст, выходим из цикла
			goto sendMetrics
		}
	}

sendMetrics:
	// Отправляем все метрики
	for _, metricRequest := range metricsToSend {
		a.sendMetric(metricRequest)
	}
}

// sendMetric отправляет одну метрику
func (a *MetricAgentService) sendMetric(metricRequest *types.UpdateMetricRequest) {
	// Отправляем метрику
	respBody, respStatus, err := a.metricFacade.SendMetric(metricRequest)
	if err != nil {
		log.Printf("Error sending metric %s:%s: %v\n", metricRequest.Type, metricRequest.Name, err)
		return
	}

	log.Printf("Sent metric %s:%s, status: %d, response: %s\n", metricRequest.Type, metricRequest.Name, respStatus, string(respBody))
}
