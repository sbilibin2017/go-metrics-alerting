package services

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"
	"os"
	"time"

	"go.uber.org/zap"
)

// Интерфейс для коллекционера метрик
type MetricsCollector interface {
	Collect() []*domain.Metrics
}

// Интерфейс для фасада отправки метрик
type MetricFacade interface {
	UpdateMetric(metric *domain.Metrics)
}

// MetricAgent структура агента для сбора и отправки метрик
type MetricAgentService struct {
	config     *configs.AgentConfig
	collectors []MetricsCollector
	facade     MetricFacade
	metricsCh  chan *domain.Metrics
	logger     *zap.Logger
}

// NewMetricAgent создает новый экземпляр MetricAgent с логгером
func NewMetricAgentService(
	config *configs.AgentConfig,
	collectors []MetricsCollector,
	facade MetricFacade,
	logger *zap.Logger,
) *MetricAgentService {
	return &MetricAgentService{
		config:     config,
		collectors: collectors,
		facade:     facade,
		metricsCh:  make(chan *domain.Metrics, 100),
		logger:     logger,
	}
}

// Run запускает агент для сбора и отправки метрик
func (a *MetricAgentService) Run(signalCh chan os.Signal) {
	a.logger.Info("Starting agent", zap.String("Address", a.config.Address), zap.Duration("PollInterval", a.config.PollInterval), zap.Duration("ReportInterval", a.config.ReportInterval))

	tickerPoll := time.NewTicker(a.config.PollInterval)
	tickerReport := time.NewTicker(a.config.ReportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	for {
		select {
		case <-signalCh:
			a.logger.Info("Received shutdown signal. Stopping agent...")
			return
		case <-tickerPoll.C:
			a.logger.Info("Collecting metrics...")
			a.collectMetrics()
		case <-tickerReport.C:
			a.logger.Info("Sending metrics...")
			a.sendMetrics()
		}
	}
}

// collectMetrics собирает метрики с использованием коллекционеров
func (a *MetricAgentService) collectMetrics() {
	for _, collector := range a.collectors {
		metrics := collector.Collect()
		for _, metric := range metrics {
			a.metricsCh <- metric
		}
	}
}

// sendMetrics отправляет метрики с помощью отправителя
func (a *MetricAgentService) sendMetrics() {
	for metric := range a.metricsCh {
		a.logger.Debug("Sending metric", zap.String("metricID", metric.ID))
		a.facade.UpdateMetric(metric)
	}
}
