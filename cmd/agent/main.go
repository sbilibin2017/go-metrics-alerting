package main

import (
	"flag"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/facades"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/strategies"

	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// main function initializes the service, loads configuration, and starts the agent
func main() {
	// Инициализация логера
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("Logger initialized successfully")

	// Инициализация конфигурации
	var config configs.AgentConfig
	parseFlagsAndEnv(&config) // передаем указатель на config

	// Инициализация фасада для отправки метрик
	facade := newFacade(&config)
	logger.Info("Facade initialized for metric reporting")

	// Создание коллектора метрик
	counterCollector, gaugeCollector := newCollectors()
	logger.Info("Metric collectors initialized", zap.Int("counterCollector", len(counterCollector.Collect())), zap.Int("gaugeCollector", len(gaugeCollector.Collect())))

	// Создаем MetricAgentService с нужными параметрами
	service := services.NewMetricAgentService(&config, facade, counterCollector, gaugeCollector)

	service.Run(logger)
}

func newLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

// parseFlagsAndEnv обрабатывает и флаги, и переменные окружения
func parseFlagsAndEnv(config *configs.AgentConfig) {
	// Парсим переменные окружения
	err := env.Parse(config)
	if err != nil {
		panic(fmt.Sprintf("Error parsing environment variables: %v", err))
	}

	// Флаги командной строки
	addressFlag := flag.String("address", "", "Address for the HTTP server")
	pollIntervalFlag := flag.Duration("poll-interval", config.PollInterval, "Interval between polls")
	reportIntervalFlag := flag.Duration("report-interval", config.ReportInterval, "Interval between reports")
	flag.Parse()

	// Если флаг пустой, то используем значения из переменных окружения
	if *addressFlag == "" && config.Address != "" {
		*addressFlag = config.Address
	}

	if *pollIntervalFlag == config.PollInterval && config.PollInterval != 0 {
		*pollIntervalFlag = config.PollInterval
	}

	if *reportIntervalFlag == config.ReportInterval && config.ReportInterval != 0 {
		*reportIntervalFlag = config.ReportInterval
	}

	// Обновляем конфигурацию с флагов
	if *addressFlag != "" {
		config.Address = *addressFlag
	}

	if *pollIntervalFlag != config.PollInterval {
		config.PollInterval = *pollIntervalFlag
	}

	if *reportIntervalFlag != config.ReportInterval {
		config.ReportInterval = *reportIntervalFlag
	}
}

func newFacade(config *configs.AgentConfig) *facades.MetricFacade {
	client := resty.New()
	return facades.NewMetricFacade(client, config)
}

func newCollectors() (*strategies.CounterMetricsCollector, *strategies.GaugeMetricsCollector) {
	return strategies.NewCounterMetricsCollector(), strategies.NewGaugeMetricsCollector()
}
