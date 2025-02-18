package agent

import (
	"flag"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/services"
	"log"
	"time" // Импортируем пакет для работы с time.Duration

	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func RunAgent() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Чтение конфигурации из переменных окружения
	var config configs.AgentConfig
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	// Обработка флага командной строки для адреса
	addressFlag := flag.String("address", "", "Address for HTTP server")
	reportIntervalFlag := flag.Int("reportInterval", 0, "Report interval in seconds")
	pollIntervalFlag := flag.Int("pollInterval", 0, "Poll interval in seconds")
	flag.Parse()

	// Приоритет значений:
	// 1. Переменная окружения (если указана)
	// 2. Флаг командной строки (если указан)
	// 3. Значение по умолчанию

	if *addressFlag != "" {
		config.Address = *addressFlag
	}

	if *reportIntervalFlag != 0 {
		// Преобразуем интервалы в time.Duration
		config.ReportInterval = time.Duration(*reportIntervalFlag) * time.Second
	}

	if *pollIntervalFlag != 0 {
		// Преобразуем интервалы в time.Duration
		config.PollInterval = time.Duration(*pollIntervalFlag) * time.Second
	}

	// Проверка значений конфигурации
	if config.Address == "" {
		log.Fatal("ADDRESS is required")
	}
	if config.ReportInterval == 0 {
		log.Fatal("REPORT_INTERVAL is required")
	}
	if config.PollInterval == 0 {
		log.Fatal("POLL_INTERVAL is required")
	}

	// Выводим конфигурацию для проверки
	fmt.Printf("Address: %s\n", config.Address)
	fmt.Printf("Report Interval: %v\n", config.ReportInterval)
	fmt.Printf("Poll Interval: %v\n", config.PollInterval)

	// Создаем HTTP клиент
	client := resty.New()

	// Запуск агента сбора метрик
	services.StartMetricAgent(&config, client)
}
