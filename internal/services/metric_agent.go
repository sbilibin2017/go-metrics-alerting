package services

import (
	"fmt"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
)

// MetricFacade - интерфейс для фасада метрик
type MetricFacade interface {
	UpdateMetric(metric types.MetricsRequest) error
}

// MetricCollector - интерфейс для сбора метрик
type MetricCollector interface {
	Collect() []types.MetricsRequest
}

// MetricAgentService - структура для сбора и отправки метрик.
type MetricAgentService struct {
	config       *configs.AgentConfig
	metricsCh    chan types.MetricsRequest
	metricFacade MetricFacade
	collectors   []MetricCollector
}

// NewMetricAgentService - создает новый сервис для сбора и отправки метрик.
func NewMetricAgentService(
	config *configs.AgentConfig,
	metricFacade MetricFacade,
	collectors []MetricCollector,
) *MetricAgentService {
	return &MetricAgentService{
		config:       config,
		metricsCh:    make(chan types.MetricsRequest, 100),
		metricFacade: metricFacade,
		collectors:   collectors,
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
			for _, collector := range s.collectors {
				metrics := collector.Collect()
				for _, metric := range metrics {
					s.metricsCh <- metric
				}
			}

		case <-tickerReport.C:
			for metric := range s.metricsCh {
				err := s.metricFacade.UpdateMetric(metric)
				if err != nil {
					fmt.Println("Error sending metric:", err)
				}
			}
		}
	}
}
