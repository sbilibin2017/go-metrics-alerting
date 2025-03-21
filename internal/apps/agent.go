package apps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"

	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultAgentAddress        = ":8080"
	DefaultAgentReportInterval = "10" // Report interval as string
	DefaultAgentPollInterval   = "2"  // Poll interval as string

	EnvAgentAddress        = "AGENT_ADDRESS"
	EnvAgentReportInterval = "AGENT_REPORT_INTERVAL"
	EnvAgentPollInterval   = "AGENT_POLL_INTERVAL"

	FlagAgentAddress        = "a"
	FlagAgentReportInterval = "r"
	FlagAgentPollInterval   = "p"

	DescriptionAgentAddress        = "Server address to report metrics"
	DescriptionAgentPollInterval   = "Interval in seconds to poll metrics"
	DescriptionAgentReportInterval = "Interval in seconds to report metrics"
)

// NewAgentCommand initializes the Cobra command for the agent configuration
func NewAgentCommand() *cobra.Command {
	var config configs.AgentConfig

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Initialize agent configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Retrieve configuration values from flags or environment variables
			config.Address = viper.GetString(FlagAgentAddress)
			if config.Address == "" {
				config.Address = viper.GetString(EnvAgentAddress)
			}
			if config.Address == "" {
				config.Address = DefaultAgentAddress
			}

			// Treat ReportInterval and PollInterval as strings
			config.ReportInterval = viper.GetString(FlagAgentReportInterval)
			if config.ReportInterval == "" {
				config.ReportInterval = viper.GetString(EnvAgentReportInterval)
			}
			if config.ReportInterval == "" {
				config.ReportInterval = DefaultAgentReportInterval
			}

			config.PollInterval = viper.GetString(FlagAgentPollInterval)
			if config.PollInterval == "" {
				config.PollInterval = viper.GetString(EnvAgentPollInterval)
			}
			if config.PollInterval == "" {
				config.PollInterval = DefaultAgentPollInterval
			}

			// Print the final configuration (you can replace this with actual usage)
			fmt.Printf("Configuration: %#v\n", config)

			// Set up the agent's context and workers
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			return runAgentApp(ctx, &config)
		},
	}

	// Set up flags for agent configuration as strings
	cmd.Flags().String(FlagAgentAddress, DefaultAgentAddress, DescriptionAgentAddress)
	cmd.Flags().String(FlagAgentReportInterval, DefaultAgentReportInterval, DescriptionAgentReportInterval)
	cmd.Flags().String(FlagAgentPollInterval, DefaultAgentPollInterval, DescriptionAgentPollInterval)

	// Bind flags to Viper
	viper.BindPFlag(FlagAgentAddress, cmd.Flags().Lookup(FlagAgentAddress))
	viper.BindPFlag(FlagAgentReportInterval, cmd.Flags().Lookup(FlagAgentReportInterval))
	viper.BindPFlag(FlagAgentPollInterval, cmd.Flags().Lookup(FlagAgentPollInterval))

	// Set up Viper to read environment variables automatically
	viper.AutomaticEnv()

	// Bind environment variables to Viper
	viper.BindEnv(FlagAgentAddress, EnvAgentAddress)
	viper.BindEnv(FlagAgentReportInterval, EnvAgentReportInterval)
	viper.BindEnv(FlagAgentPollInterval, EnvAgentPollInterval)

	return cmd
}

// StartMetricAgentWorker запускает сбор и отправку метрик с graceful shutdown
func runAgentApp(ctx context.Context, config *configs.AgentConfig) error {
	var metrics []types.Metrics
	var pollCount int64

	// Parse PollInterval and ReportInterval from config
	pollInterval, err := strconv.Atoi(config.PollInterval)
	if err != nil {
		return err
	}

	reportInterval, err := strconv.Atoi(config.ReportInterval)
	if err != nil {
		return err
	}

	// Setup tickers for polling and reporting intervals
	tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-tickerPoll.C:
			collectMetrics(&metrics, &pollCount)
		case <-tickerReport.C:
			// Handling retriable errors for reportMetrics
			reportMetrics(config, metrics)
			resetMetrics(&metrics, &pollCount)
		case <-sigChan:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}

// collectMetrics собирает метрики системы
func collectMetrics(metrics *[]types.Metrics, pollCount *int64) {
	newFloat64 := func(v float64) *float64 { return &v }
	newInt64 := func(v int64) *int64 { return &v }

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Collect memory stats metrics
	*metrics = []types.Metrics{
		{ID: "Alloc", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.Alloc))},
		{ID: "BuckHashSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.BuckHashSys))},
		{ID: "Frees", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.Frees))},
		{ID: "GCCPUFraction", Type: string(types.Gauge), Delta: nil, Value: newFloat64(memStats.GCCPUFraction)},
		{ID: "GCSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.GCSys))},
		{ID: "HeapAlloc", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.HeapAlloc))},
		{ID: "HeapIdle", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.HeapIdle))},
		{ID: "HeapInuse", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.HeapInuse))},
		{ID: "HeapObjects", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.HeapObjects))},
		{ID: "HeapReleased", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.HeapReleased))},
		{ID: "HeapSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.HeapSys))},
		{ID: "LastGC", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.LastGC))},
		{ID: "Lookups", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.Lookups))},
		{ID: "MCacheInuse", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.MCacheInuse))},
		{ID: "MCacheSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.MCacheSys))},
		{ID: "MSpanInuse", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.MSpanInuse))},
		{ID: "MSpanSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.MSpanSys))},
		{ID: "Mallocs", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.Mallocs))},
		{ID: "NextGC", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.NextGC))},
		{ID: "NumForcedGC", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.NumForcedGC))},
		{ID: "NumGC", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.NumGC))},
		{ID: "OtherSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.OtherSys))},
		{ID: "PauseTotalNs", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.PauseTotalNs))},
		{ID: "StackInuse", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.StackInuse))},
		{ID: "StackSys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.StackSys))},
		{ID: "Sys", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.Sys))},
		{ID: "TotalAlloc", Type: string(types.Gauge), Delta: nil, Value: newFloat64(float64(memStats.TotalAlloc))},
		{ID: "RandomValue", Type: string(types.Gauge), Delta: nil, Value: newFloat64(rand.Float64())},
	}

	// Increment PollCount
	*pollCount += 1
	*metrics = append(*metrics,
		types.Metrics{ID: "PollCount", Type: string(types.Counter), Delta: newInt64(*pollCount), Value: nil},
	)
}

// resetMetrics сбрасывает собранные метрики и счетчик PollCount
func resetMetrics(metrics *[]types.Metrics, pollCount *int64) {
	*metrics = []types.Metrics{}
	*pollCount = 0
}

// reportMetrics отправляет метрики на сервер с обработкой retriable-ошибок
func reportMetrics(config *configs.AgentConfig, metrics []types.Metrics) {
	// Параметры для retry
	maxRetries := 3
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	// Marshal metrics to JSON
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return
	}

	// Prepare the URL and check if it has a scheme
	url := config.Address + "/updates/"
	if !strings.HasPrefix(config.Address, "http://") && !strings.HasPrefix(config.Address, "https://") {
		url = "http://" + url
	}

	// Retry sending the metrics
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err == nil {
			// Successful request
			defer resp.Body.Close()

			// Reset metrics after reporting
			return
		}

		// Log error and retry if attempts left
		if attempt < maxRetries-1 {
			time.Sleep(retryIntervals[attempt])
		} else {
			return
		}
	}

	return
}
