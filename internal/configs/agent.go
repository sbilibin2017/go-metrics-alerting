package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func LoadAgentConfigFromEnv() (*AgentConfig, error) {
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	address := os.Getenv("ADDRESS")
	reportIntervalDuration, _ := time.ParseDuration(os.Getenv("REPORT_INTERVAL"))
	pollIntervalDuration, _ := time.ParseDuration(os.Getenv("POLL_INTERVAL"))
	return &AgentConfig{
		Address:        address,
		ReportInterval: reportIntervalDuration,
		PollInterval:   pollIntervalDuration,
	}, nil
}

func LoadAgentConfigFromFlags() (*AgentConfig, error) {
	var config AgentConfig

	// Создаем новый корневой командный объект
	var rootCmd = &cobra.Command{
		Use:   "agent",
		Short: "Agent to collect and report metrics",
		Run: func(cmd *cobra.Command, args []string) {
			// Не нужно выполнять здесь, можно оставить для обработки флагов
		},
	}

	// Добавляем флаги для конфигурации
	rootCmd.Flags().StringVarP(&config.Address, "address", "a", "localhost:8080", "HTTP server endpoint")
	rootCmd.Flags().DurationVarP(&config.ReportInterval, "reportInterval", "r", 10*time.Second, "Report interval (default 10s)")
	rootCmd.Flags().DurationVarP(&config.PollInterval, "pollInterval", "p", 2*time.Second, "Poll interval (default 2s)")

	// Выполняем команду и парсим флаги
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	// Возвращаем структуру с конфигурацией
	return &config, nil
}
