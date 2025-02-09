package configs

import (
	"time"
)

// Структура конфигурации для агент-сервиса
type AgentConfig struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // Интервал опроса
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // Интервал отчётов
	Address        string        `env:"ADDRESS"`         // Адрес агента
}
