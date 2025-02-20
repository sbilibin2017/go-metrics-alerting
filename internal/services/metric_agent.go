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
	Collect() []*domain.Metric
}

// Интерфейс для фасада отправки метрик
type MetricFacade interface {
	UpdateMetric(metric *domain.Metric)
}

// MetricAgent структура агента для сбора и отправки метрик
type MetricAgent struct {
	config     *configs.AgentConfig
	collectors []MetricsCollector
	facade     MetricFacade
	logger     *zap.Logger
}

// NewMetricAgent создает новый экземпляр MetricAgent с логгером
func NewMetricAgent(config *configs.AgentConfig, collectors []MetricsCollector, facade MetricFacade, logger *zap.Logger) *MetricAgent {
	return &MetricAgent{
		config:     config,
		collectors: collectors,
		facade:     facade,
		logger:     logger,
	}
}

// Run запускает агент для сбора и отправки метрик
func (a *MetricAgent) Run(signalCh chan os.Signal) {
	a.logger.Info("Starting agent", zap.String("Address", a.config.Address), zap.Duration("PollInterval", a.config.PollInterval), zap.Duration("ReportInterval", a.config.ReportInterval))

	metricsCh := make(chan *domain.Metric, 100)

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
			a.collectMetrics(metricsCh)
		case <-tickerReport.C:
			a.logger.Info("Sending metrics...")
			a.sendMetrics(metricsCh)
		}
	}
}

// collectMetrics собирает метрики с использованием коллекционеров
func (a *MetricAgent) collectMetrics(metricsCh chan *domain.Metric) {
	for _, collector := range a.collectors {
		metrics := collector.Collect()
		for _, metric := range metrics {
			metricsCh <- metric
		}
	}
}

// sendMetrics отправляет метрики с помощью отправителя
func (a *MetricAgent) sendMetrics(metricsCh chan *domain.Metric) {
	for metric := range metricsCh {
		a.logger.Debug("Sending metric", zap.String("metricID", metric.ID), zap.String("metricValue", metric.Value))
		a.facade.UpdateMetric(metric)
	}
}
