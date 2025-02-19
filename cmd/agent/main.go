package main

import (
	"flag"
	"go-metrics-alerting/internal/agent"
	"go-metrics-alerting/internal/configs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
)

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	var config configs.AgentConfig

	env.Parse(&config)

	// Печать финальных значений конфигурации после обработки флагов
	log.Printf("Final config: Address=%s, ReportInterval=%v, PollInterval=%v", config.Address, config.ReportInterval, config.PollInterval)

	addressFlag := flag.String("address", "", "Address for HTTP server")
	reportIntervalFlag := flag.Int("reportInterval", 0, "Report interval in seconds")
	pollIntervalFlag := flag.Int("pollInterval", 0, "Poll interval in seconds")
	flag.Parse()

	if config.Address == "" {
		config.Address = *addressFlag
	}

	if config.ReportInterval == 0 {
		config.ReportInterval = time.Duration(*reportIntervalFlag) * time.Second
	}

	if config.PollInterval == 0 {
		config.PollInterval = time.Duration(*pollIntervalFlag) * time.Second
	}

	client := resty.New()

	agent.StartAgent(signalCh, &config, client)

}
