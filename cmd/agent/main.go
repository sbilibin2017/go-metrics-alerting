package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/logger"

	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
)

func main() {
	config := &configs.AgentConfig{}
	if err := env.Parse(config); err != nil {
		logger.Logger.Fatalf("Error parsing environment variables: %v", err)
	}

	if config.Address == "" {
		config.Address = ":8080"
	}
	if config.PollInterval == 0 {
		config.PollInterval = 2 * time.Second
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = 10 * time.Second
	}

	address := flag.String("address", config.Address, "Address of the agent (default: :8080)")
	reportInterval := flag.Int("report-interval", int(config.ReportInterval.Seconds()), "Report interval in seconds")
	pollInterval := flag.Int("poll-interval", int(config.PollInterval.Seconds()), "Poll interval in seconds")

	flag.Parse()

	if *address != config.Address {
		config.Address = *address
	}
	if *reportInterval != int(config.ReportInterval.Seconds()) {
		config.ReportInterval = time.Duration(*reportInterval) * time.Second
	}
	if *pollInterval != int(config.PollInterval.Seconds()) {
		config.PollInterval = time.Duration(*pollInterval) * time.Second
	}

	client := resty.New()

	agentService := &services.MetricAgentService{
		APIClient:      client,
		PollInterval:   config.PollInterval,
		ReportInterval: config.ReportInterval,
		MetricChannel:  make(chan types.UpdateMetricValueRequest, 100),
		Shutdown:       make(chan os.Signal, 1),
		Address:        config.Address,
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go agentService.Start()

	sig := <-signalChannel
	agentService.Shutdown <- sig
}
