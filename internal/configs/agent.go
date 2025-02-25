package configs

import "time"

type AgentConfig struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`    // Интервал опроса, дефолтное значение 2 секунды
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"` // Интервал отчётов, дефолтное значение 10 секунд
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
}
