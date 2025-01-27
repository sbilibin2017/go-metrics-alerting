package main

import (
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/facades"
	"go-metrics-alerting/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Получаем конфигурацию агента через флаги и переменные окружения
	agentConfig := configs.NewAgentConfig()

	// Создаем экземпляр API клиента
	apiClient := engines.NewApiClient()

	// Создаем фасад для отправки метрик, передавая объект agentConfig (который реализует интерфейс)
	metricFacade := facades.NewMetricFacade(apiClient, agentConfig)

	// Создаем сервис для агентов
	metricAgentService := services.NewMetricAgentService(
		agentConfig,
		metricFacade,
	)

	// Запуск агента
	go metricAgentService.Start()

	// Необходимо, чтобы основной процесс не завершался
	select {}
}
