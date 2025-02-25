package services

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// MetricFacade - интерфейс для фасада, который обновляет метрики
type MetricFacade interface {
	UpdateMetric(metric *types.UpdateMetricBodyRequest)
}

// MetricsCollector - интерфейс для коллектора метрик
type MetricsCollector interface {
	Collect() []*types.UpdateMetricBodyRequest
}

// MetricAgentService - структура для сервиса агента метрик
type MetricAgentService struct {
	config           *configs.AgentConfig
	facade           MetricFacade
	counterCollector MetricsCollector
	gaugeCollector   MetricsCollector
	metricsToReport  []*types.UpdateMetricBodyRequest
}

// NewMetricAgentService - инициализация нового агента метрик
func NewMetricAgentService(config *configs.AgentConfig, facade MetricFacade, counterCollector, gaugeCollector MetricsCollector) *MetricAgentService {
	return &MetricAgentService{
		config:           config,
		facade:           facade,
		counterCollector: counterCollector,
		gaugeCollector:   gaugeCollector,
		metricsToReport:  []*types.UpdateMetricBodyRequest{},
	}
}

// Run - основной цикл агента метрик
func (service *MetricAgentService) Run(logger *zap.Logger) {
	logger.Info("Metric agent service is starting...")

	// Создание тикеров с интервалами из конфигурации
	pollTicker := time.NewTicker(service.config.PollInterval)
	reportTicker := time.NewTicker(service.config.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	// Канал для ловли сигналов от операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	for {
		select {
		case <-pollTicker.C:
			service.pollMetrics(logger)
		case <-reportTicker.C:
			service.reportMetrics(logger)
		case <-sigChan:
			return
		}
	}
}

// pollMetrics - сбор метрик
func (service *MetricAgentService) pollMetrics(logger *zap.Logger) {
	logger.Debug("Polling for metrics...")

	// Собираем метрики
	gaugeMetrics := service.gaugeCollector.Collect()
	counterMetrics := service.counterCollector.Collect()

	// Пропускаем, если метрики не собраны
	if len(gaugeMetrics) == 0 && len(counterMetrics) == 0 {
		logger.Debug("No metrics to collect, skipping polling cycle")
		return
	}

	// Добавляем собранные метрики в массив
	service.metricsToReport = append(service.metricsToReport, gaugeMetrics...)
	service.metricsToReport = append(service.metricsToReport, counterMetrics...)

	logger.Debug("Collected metrics", zap.Int("gaugeMetrics", len(gaugeMetrics)), zap.Int("counterMetrics", len(counterMetrics)))
}

// reportMetrics - отправка метрик на сервер
func (service *MetricAgentService) reportMetrics(logger *zap.Logger) {
	if len(service.metricsToReport) > 0 {
		logger.Debug("Reporting metrics...")
		for _, metric := range service.metricsToReport {
			service.facade.UpdateMetric(metric)
		}
		service.metricsToReport = nil
		logger.Debug("All metrics have been reported, clearing metrics array.")
	} else {
		logger.Debug("No metrics to report, skipping reporting cycle")
	}
}
