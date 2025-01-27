package configs

import (
	"flag"
	"fmt"
	"time"
)

// Интерфейс для конфигурации агента
type AgentConfigInterface interface {
	GetPollInterval() time.Duration
	GetReportInterval() time.Duration
	GetServerURL() string
}

// Структура конфигурации для агента
type AgentConfig struct {
	PollInterval   time.Duration // Интервал для опроса метрик
	ReportInterval time.Duration // Интервал для отправки метрик
	ServerURL      string        // URL сервера для отправки метрик
}

// Функция для создания конфигурации агента с флагами и переменными окружения
func NewAgentConfig() *AgentConfig {
	// Флаги для конфигурации агента
	serverURL := flag.String("a", "", "Адрес эндпоинта HTTP-сервера (например, http://localhost:8080)")
	reportInterval := flag.Int("r", 10, "Частота отправки метрик на сервер (в секундах)")
	pollInterval := flag.Int("p", 2, "Частота опроса метрик (в секундах)")

	flag.Parse()

	// Получаем значение из переменной окружения для серверного URL, если оно установлено
	finalServerURL := GetEnvOrDefault(*serverURL, "SERVER_URL", "http://localhost:8080")

	// Получаем значения из переменных окружения для интервалов
	pollIntervalEnv := GetEnvOrDefault(fmt.Sprint(*pollInterval), "POLL_INTERVAL", "2")
	reportIntervalEnv := GetEnvOrDefault(fmt.Sprint(*reportInterval), "REPORT_INTERVAL", "10")

	// Преобразуем интервалы в тип time.Duration
	pollIntervalDuration, err := time.ParseDuration(pollIntervalEnv + "s")
	if err != nil {
		pollIntervalDuration = 2 * time.Second
	}

	reportIntervalDuration, err := time.ParseDuration(reportIntervalEnv + "s")
	if err != nil {
		reportIntervalDuration = 10 * time.Second
	}

	return &AgentConfig{
		PollInterval:   pollIntervalDuration,
		ReportInterval: reportIntervalDuration,
		ServerURL:      finalServerURL,
	}
}

// Реализация методов интерфейса AgentConfigInterface для AgentConfig
func (config *AgentConfig) GetPollInterval() time.Duration {
	return config.PollInterval
}

func (config *AgentConfig) GetReportInterval() time.Duration {
	return config.ReportInterval
}

func (config *AgentConfig) GetServerURL() string {
	return config.ServerURL
}
