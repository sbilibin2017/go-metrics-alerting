package configs

type AgentConfig struct {
	Address        string
	ReportInterval string
	PollInterval   string
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{}
}
