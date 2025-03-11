package configs

import "time"

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}
