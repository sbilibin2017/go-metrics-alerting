package main

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/facades"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/strategies"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func main() {
	// Считываем конфигурацию из переменных окружения
	var config configs.AgentConfig
	env.Parse(&config)

	// Парсинг флагов командной строки
	addressFlag := flag.String("address", "", "Address for the HTTP server")
	pollIntervalFlag := flag.Duration("poll-interval", config.PollInterval, "Interval between polls")
	reportIntervalFlag := flag.Duration("report-interval", config.ReportInterval, "Interval between reports")
	flag.Parse()

	// Переопределение значений, если заданы флаги
	if *addressFlag != "" {
		config.Address = *addressFlag
	}
	if *pollIntervalFlag != config.PollInterval {
		config.PollInterval = *pollIntervalFlag
	}
	if *reportIntervalFlag != config.ReportInterval {
		config.ReportInterval = *reportIntervalFlag
	}

	// Создаем логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Создаем клиент для отправки запросов
	client := resty.New()

	// Инициализируем стратегии для сбора метрик
	collectorCounterStrategy := strategies.NewCounterMetricsCollector()
	collectorGaugeStrategy := strategies.NewGaugeMetricsCollector()

	// Создаем фасад для метрик
	facade := facades.NewMetricFacade(
		client,
		&config,
		logger,
	)

	// Создаем сервис агента
	metricAgentService := services.NewMetricAgentService(
		&config,
		collectorCounterStrategy,
		collectorGaugeStrategy,
		facade,
		logger,
	)

	// Создаем канал для перехвата сигналов завершения
	signalCh := make(chan os.Signal, 1)

	// Регистрируем сигнал завершения (например, SIGINT или SIGTERM)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Запуск агента в горутине
	go metricAgentService.Run(signalCh)

	// Ожидаем сигнал завершения
	<-signalCh
	logger.Info("Shutting down the agent gracefully...")
}
