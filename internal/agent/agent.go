package agent

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func Run() {
	// Чтение конфигурации из окружения
	config := &configs.AgentConfig{}

	godotenv.Load()
	env.Parse(config)

	// Значения по умолчанию
	if config.Address == "" {
		config.Address = ":8080"
	}
	if config.PollInterval == 0 {
		config.PollInterval = 2 * time.Second
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = 10 * time.Second
	}

	// Обработка флагов
	address := flag.String("address", config.Address, "Address of the agent (default: :8080)")
	reportInterval := flag.Int("report-interval", int(config.ReportInterval.Seconds()), "Report interval in seconds")
	pollInterval := flag.Int("poll-interval", int(config.PollInterval.Seconds()), "Poll interval in seconds")
	flag.Parse()

	// Обновление конфигурации, если флаги изменены
	if *address != config.Address {
		config.Address = *address
	}
	if *reportInterval != int(config.ReportInterval.Seconds()) {
		config.ReportInterval = time.Duration(*reportInterval) * time.Second
	}
	if *pollInterval != int(config.PollInterval.Seconds()) {
		config.PollInterval = time.Duration(*pollInterval) * time.Second
	}

	// Создание объекта MetricAgentService
	agentService := services.NewMetricAgentService(config)

	// Канал для сигнала остановки
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервиса сбора метрик в горутине
	go agentService.Start()

	// Ожидание сигнала для завершения работы
	sig := <-signalChannel

	println("Received signal:", sig)

}
