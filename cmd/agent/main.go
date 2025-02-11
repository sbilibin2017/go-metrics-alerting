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
	// Load configuration from environment variables
	config := &configs.AgentConfig{}
	if err := env.Parse(config); err != nil {
		logger.Logger.Fatalf("Error parsing environment variables: %v", err)
	}

	// Set default values for the configuration if not set
	if config.Address == "" {
		config.Address = ":8080"
	}
	if config.PollInterval == 0 {
		config.PollInterval = 2 * time.Second
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = 10 * time.Second
	}

	// Read command-line flags (they have second priority)
	address := flag.String("address", config.Address, "Address of the agent (default: :8080)")
	reportInterval := flag.Int("report-interval", int(config.ReportInterval.Seconds()), "Report interval in seconds")
	pollInterval := flag.Int("poll-interval", int(config.PollInterval.Seconds()), "Poll interval in seconds")

	// Parse the flags
	flag.Parse()

	// Override the configuration with command-line flags if provided
	if *address != config.Address {
		config.Address = *address
	}
	if *reportInterval != int(config.ReportInterval.Seconds()) {
		config.ReportInterval = time.Duration(*reportInterval) * time.Second
	}
	if *pollInterval != int(config.PollInterval.Seconds()) {
		config.PollInterval = time.Duration(*pollInterval) * time.Second
	}

	// Log the configuration and the service start-up
	logger.Logger.Infof("Starting Metric Agent Service on %s with Poll Interval: %s and Report Interval: %s", config.Address, config.PollInterval, config.ReportInterval)

	// Create a new resty client
	client := resty.New()

	// Instantiate the MetricAgentService with the Resty client, the pollInterval, and the reportInterval
	agentService := &services.MetricAgentService{
		APIClient:      client,
		PollInterval:   config.PollInterval,
		ReportInterval: config.ReportInterval,
		MetricChannel:  make(chan types.UpdateMetricValueRequest, 100),
		Shutdown:       make(chan os.Signal, 1),
		Address:        config.Address, // Or another base URL if needed
	}

	// Create a channel to catch termination signals (e.g. Ctrl+C)
	signalChannel := make(chan os.Signal, 1)
	// Register for SIGINT (Ctrl+C) and SIGTERM (process termination)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Start the agent service in a new goroutine
	go agentService.Start()

	// Block the main goroutine, waiting for termination signal
	sig := <-signalChannel
	logger.Logger.Infof("Received signal: %s. Shutting down agent...", sig)

	// Gracefully shut down the agent service (you may want to add a Stop() method to agentService)
	// For example, you can use `agentService.Stop()` if such method exists
	agentService.Shutdown <- sig // Or call a shutdown function on the service, if it exists

	// Ensure that the agent service cleans up properly before the program exits
	// You can also wait for the agent to complete any active operations or goroutines
	// if necessary before shutting down completely.

	// Log the shutdown process
	logger.Logger.Info("Agent service shutdown complete.")
}
