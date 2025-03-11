package apps

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/apps/helpers"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"
	"net/http"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/spf13/cobra"
)

// Дефолтные значения
const (
	FlagAgentAddress   = "address"
	FlagReportInterval = "report-interval"
	FlagPollInterval   = "poll-interval"

	FlagAgentAddressShort   = "a"
	FlagReportIntervalShort = "r"
	FlagPollIntervalShort   = "p"

	EnvAgentAddress   = "ADDRESS"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvPollInterval   = "POLL_INTERVAL"

	DescAgentAddress   = "Server address"
	DescReportInterval = "Metrics report interval"
	DescPollInterval   = "Metrics poll interval"

	DefaultAgentAddress   = "localhost:8080"
	DefaultReportInterval = 10 * time.Second
	DefaultPollInterval   = 2 * time.Second

	UseAgent   = "agent"
	ShortAgent = "Start the agent with specified parameters"
)

func NewAgentAppCommand() *cobra.Command {
	var (
		address        string
		reportInterval time.Duration
		pollInterval   time.Duration
	)

	// Используем константы для загрузки конфигурации из переменных окружения
	addressEnv := helpers.GetStringFromEnv(EnvAgentAddress)
	if addressEnv != nil {
		address = *addressEnv
	} else {
		address = DefaultAgentAddress
	}

	reportIntervalEnv := helpers.GetDurationFromEnv(EnvReportInterval)
	if reportIntervalEnv != nil {
		reportInterval = *reportIntervalEnv
	} else {
		reportInterval = DefaultReportInterval
	}

	pollIntervalEnv := helpers.GetDurationFromEnv(EnvPollInterval)
	if pollIntervalEnv != nil {
		pollInterval = *pollIntervalEnv
	} else {
		pollInterval = DefaultPollInterval
	}

	cmd := &cobra.Command{
		Use:   UseAgent,
		Short: ShortAgent,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := configs.AgentConfig{
				Address:        address,
				ReportInterval: reportInterval,
				PollInterval:   pollInterval,
			}
			runAgent(context.Background(), &cfg)
		},
	}

	// Используем константы для флагов агента
	cmd.Flags().StringVarP(
		&address,
		FlagAgentAddress,
		FlagAgentAddressShort,
		address,
		DescAgentAddress,
	)

	cmd.Flags().DurationVarP(
		&reportInterval,
		FlagReportInterval,
		FlagReportIntervalShort,
		reportInterval,
		DescReportInterval,
	)

	cmd.Flags().DurationVarP(
		&pollInterval,
		FlagPollInterval,
		FlagPollIntervalShort,
		pollInterval,
		DescPollInterval,
	)

	return cmd
}

var metrics []types.Metrics

// Start запускает сбор метрик и их отправку на сервер.
func runAgent(ctx context.Context, config *configs.AgentConfig) {
	tickerPoll := time.NewTicker(config.PollInterval)
	tickerReport := time.NewTicker(config.ReportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tickerPoll.C:
			metrics = append(metrics, collectMetrics()...)
		case <-tickerReport.C:
			reportMetrics(ctx, config, metrics)
			metrics = nil
		}
	}
}

// CollectMetrics собирает метрики из пакета runtime.
func collectMetrics() []types.Metrics {
	var metrics []types.Metrics

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	updateCounter := func() func() int64 {
		var counter int64
		return func() int64 {
			counter += 1
			return counter
		}
	}()

	// ptrToFloat64 помогает создать указатель на float64.
	ptrToFloat64 := func(v float64) *float64 {
		return &v
	}

	// ptrToInt64 помогает создать указатель на int64.
	ptrToInt64 := func(v int64) *int64 {
		return &v
	}

	// Сбор метрик типа gauge (float64)
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "Alloc", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.Alloc)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "BuckHashSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.BuckHashSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "Frees", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.Frees)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "GCCPUFraction", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.GCCPUFraction)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "GCSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.GCSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "HeapAlloc", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.HeapAlloc)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "HeapIdle", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.HeapIdle)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "HeapInuse", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.HeapInuse)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "HeapObjects", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.HeapObjects)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "HeapReleased", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.HeapReleased)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "HeapSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.HeapSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "LastGC", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.LastGC)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "Lookups", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.Lookups)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "MCacheInuse", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.MCacheInuse)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "MCacheSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.MCacheSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "MSpanInuse", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.MSpanInuse)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "MSpanSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.MSpanSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "Mallocs", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.Mallocs)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "NextGC", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.NextGC)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "NumForcedGC", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.NumForcedGC)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "NumGC", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.NumGC)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "OtherSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.OtherSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "PauseTotalNs", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.PauseTotalNs)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "StackInuse", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.StackInuse)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "StackSys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.StackSys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "Sys", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.Sys)),
	})
	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "TotalAlloc", Type: types.Gauge},
		Value:    ptrToFloat64(float64(ms.TotalAlloc)),
	})

	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "RandomValue", Type: types.Gauge},
		Value:    ptrToFloat64(float64(time.Now().UnixNano() % 100)),
	})

	metrics = append(metrics, types.Metrics{
		MetricID: types.MetricID{ID: "PollCount", Type: types.Counter},
		Delta:    ptrToInt64(updateCounter()),
	})

	return metrics
}

// ReportMetrics отправляет метрики на сервер с логикой повторных попыток.
func reportMetrics(ctx context.Context, config *configs.AgentConfig, metrics []types.Metrics) error {
	client := resty.New()
	url := config.Address + "/update/"

	// Определение интервалов повторных попыток
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	var lastErr error
	for _, wait := range retryIntervals {
		resp, err := client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetBody(metrics).
			Post(url)

		if err == nil && resp.StatusCode() == http.StatusOK {
			return nil // Успешная отправка
		}

		// Если ошибка временная — ждём и пробуем снова
		if err != nil && isRetriableError(err) {
			lastErr = err
			time.Sleep(wait)
			continue
		}

		// Если ошибка не временная, прерываем попытки
		if err != nil {
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			return errors.New("failed to report metrics, status: " + resp.Status())
		}
	}

	return lastErr
}

// isRetriableError проверяет, является ли ошибка временной и требует ли повторной попытки.
func isRetriableError(err error) bool {
	if err == nil {
		return false
	}

	// Проверка сетевых ошибок Resty
	var urlErr *resty.Response
	if errors.As(err, &urlErr) {
		if urlErr.StatusCode() == http.StatusServiceUnavailable || urlErr.StatusCode() == http.StatusGatewayTimeout {
			return true
		}
	}

	// Проверка PostgreSQL ошибок
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case pgerrcode.ConnectionException,
			pgerrcode.ConnectionDoesNotExist,
			pgerrcode.ConnectionFailure,
			pgerrcode.TransactionResolutionUnknown:
			return true
		}
	}

	return false
}
