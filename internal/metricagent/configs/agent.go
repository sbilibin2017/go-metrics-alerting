package configs

import (
	"time"

	"github.com/caarlos0/env/v6"
)

// Структура конфигурации для агент-сервиса
type AgentConfig struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // Интервал опроса
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // Интервал отчётов
	Address        string        `env:"ADDRESS"`         // Адрес агента
}

// Функция для загрузки конфигурации
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Устанавливаем значения по умолчанию, если не заданы
	if cfg.Address == "" {
		cfg.Address = ":8080"
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 2 * time.Second
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = 10 * time.Second
	}

	return cfg, nil
}
