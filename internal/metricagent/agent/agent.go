package agent

import (
	"flag"
	"go-metrics-alerting/internal/apiclient"
	"go-metrics-alerting/internal/apiclient/facades"
	"go-metrics-alerting/internal/metricagent/configs"
	"go-metrics-alerting/internal/metricagent/services"
	"go-metrics-alerting/pkg/logger"
	"time"
)

// Функция для запуска агента
func Run() {
	// Загружаем конфигурацию из переменных окружения
	config, err := configs.LoadAgentConfig()
	if err != nil {
		logger.Logger.Fatalf("Failed to load agent config: %v", err)
	}

	// Считываем флаги командной строки (флаги имеют второй приоритет)
	address := flag.String("address", config.Address, "Address of the agent (default: :8080)")
	reportInterval := flag.Int("report-interval", int(config.ReportInterval.Seconds()), "Report interval in seconds")
	pollInterval := flag.Int("poll-interval", int(config.PollInterval.Seconds()), "Poll interval in seconds")

	// Парсим флаги
	flag.Parse()

	// Переопределяем параметры конфигурации, если указаны флаги
	if *address != config.Address {
		config.Address = *address
	}
	if *reportInterval != int(config.ReportInterval.Seconds()) {
		config.ReportInterval = time.Duration(*reportInterval) * time.Second
	}
	if *pollInterval != int(config.PollInterval.Seconds()) {
		config.PollInterval = time.Duration(*pollInterval) * time.Second
	}

	// Логирование запуска агента
	logger.Logger.Infof("Starting Metric Agent Service on %s with Poll Interval: %s and Report Interval: %s", config.Address, config.PollInterval, config.ReportInterval)

	// Создаем клиент API
	apiClient := apiclient.NewRestyClient()

	// Создаем фасад для работы с метриками
	metricFacade := facades.NewMetricFacade(apiClient, config.Address)

	// Создаем сервис для сбора и отправки метрик
	agentService := services.NewMetricAgentService(metricFacade, config.PollInterval, config.ReportInterval)

	// Запускаем сервис
	agentService.Start()
}
