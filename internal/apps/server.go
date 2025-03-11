package apps

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/apis/handlers"
	"go-metrics-alerting/internal/apis/responses"
	"go-metrics-alerting/internal/apis/routers"
	"go-metrics-alerting/internal/apps/helpers"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/engines/db"
	"go-metrics-alerting/internal/engines/file"
	"go-metrics-alerting/internal/engines/server"
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
)

const (
	// Флаги командной строки
	FlagAddress        = "address"
	FlagReportInterval = "report-interval"
	FlagPollInterval   = "poll-interval"

	// Короткие флаги
	FlagAddressShort        = "a"
	FlagReportIntervalShort = "r"
	FlagPollIntervalShort   = "p"

	// Переменные окружения
	EnvAddress        = "ADDRESS"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvPollInterval   = "POLL_INTERVAL"

	// Описания флагов
	DescAddress        = "Server address"
	DescReportInterval = "Metrics report interval"
	DescPollInterval   = "Metrics poll interval"

	// Дефолтные значения
	DefaultAddress        = "localhost:8080"
	DefaultReportInterval = 10 * time.Second
	DefaultPollInterval   = 2 * time.Second

	Use   = "agent"
	Short = "Start the agent with specified parameters"
)

func NewServerAppCommand() *cobra.Command {
	var (
		address         string
		databaseDSN     string
		fileStoragePath string
		storeInterval   time.Duration
		restore         bool
	)

	// Load configuration from environment variables or flags
	addressEnv := helpers.GetStringFromEnv(EnvAddress)
	if addressEnv != nil {
		address = *addressEnv
	} else {
		address = DefaultAddress
	}

	databaseDSNEnv := helpers.GetStringFromEnv(EnvDatabaseDSN)
	if databaseDSNEnv != nil {
		databaseDSN = *databaseDSNEnv
	} else {
		databaseDSN = DefaultDatabaseDSN
	}

	fileStoragePathEnv := helpers.GetStringFromEnv(EnvFileStoragePath)
	if fileStoragePathEnv != nil {
		fileStoragePath = *fileStoragePathEnv
	} else {
		fileStoragePath = DefaultFileStoragePath
	}

	storeIntervalEnv := helpers.GetDurationFromEnv(EnvStoreInterval)
	if storeIntervalEnv != nil {
		storeInterval = *storeIntervalEnv
	} else {
		storeInterval = DefaultStoreInterval
	}

	restoreEnv := helpers.GetBoolFromEnv(EnvRestore)
	if restoreEnv != nil {
		restore = *restoreEnv
	} else {
		restore = DefaultRestore
	}

	cmd := &cobra.Command{
		Use:   Use,
		Short: Short,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := configs.NewServerConfig(address, databaseDSN, fileStoragePath, storeInterval, restore)
			runServer(context.Background(), cfg)
		},
	}

	cmd.Flags().StringVarP(
		&address,
		FlagAddress,
		FlagAddressShort,
		address,
		DescAddress,
	)

	cmd.Flags().StringVarP(
		&databaseDSN,
		FlagDatabaseDSN,
		FlagDatabaseDSNShort,
		databaseDSN,
		DescDatabaseDSN,
	)

	cmd.Flags().StringVarP(
		&fileStoragePath,
		FlagFileStoragePath,
		FlagFileStoragePathShort,
		fileStoragePath,
		DescFileStoragePath,
	)

	cmd.Flags().DurationVarP(
		&storeInterval,
		FlagStoreInterval,
		FlagStoreIntervalShort,
		storeInterval,
		DescStoreInterval,
	)

	cmd.Flags().BoolVarP(
		&restore,
		FlagRestore,
		FlagRestoreShort,
		restore,
		DescRestore,
	)

	return cmd
}

func runServer(ctx context.Context, config *configs.ServerConfig) {
	var d *db.DB
	if config.DatabaseDSN != "" {
		d = db.NewDB()
	}

	var f *file.File
	if config.FileStoragePath != "" {
		f = file.NewFile()
	}

	// Initialize repositories and services
	mainRepo, fileRepo := repositories.NewMetricRepository(d, f)
	svc := services.NewMetricService(mainRepo)
	h := handlers.NewMetricHandler(svc)
	hr := NewHealthRouter(d)
	r := routers.NewMetricRouter(h)
	srv := server.NewServer(config.Address)
	srv.AddRouter(r, "/")
	srv.AddRouter(hr, "/ping")

	// Load metrics from file before starting the server
	loadMetricsFromFile(ctx, config, mainRepo, fileRepo)

	// Start saving metrics to file in a separate goroutine
	go saveMetricsToFile(ctx, config, mainRepo, fileRepo)

	srv.Run(ctx)
}

// MetricsRepository
type MetricsRepository interface {
	SaveMetrics(ctx context.Context, metrics []*types.Metrics) bool
	ListMetrics(ctx context.Context) []*types.Metrics
}

func loadMetricsFromFile(
	ctx context.Context,
	config *configs.ServerConfig,
	mainRepo MetricsRepository,
	fileRepo MetricsRepository,
) {
	if config.FileStoragePath == "" {
		return
	}
	metrics := fileRepo.ListMetrics(ctx)
	mainRepo.SaveMetrics(ctx, metrics)
}

func saveMetricsToFile(
	ctx context.Context,
	config *configs.ServerConfig,
	mainRepo MetricsRepository,
	fileRepo MetricsRepository,
) {
	if config.DatabaseDSN == "" {
		return
	}
	ticker := time.NewTicker(config.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := mainRepo.ListMetrics(ctx)
			fileRepo.SaveMetrics(ctx, metrics)
		}
	}
}

// NewHealthRouter creates a new router with the /ping healthcheck route.
func NewHealthRouter(db *db.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", PingHandler(db))
	return r
}

// PingHandler checks the database connection and returns the appropriate HTTP status
func PingHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			responses.InternalServerErrorResponse(w, errors.New("db is not initialized"))
			return
		}
		err := db.Ping()
		if err != nil {
			responses.InternalServerErrorResponse(w, err)
			return
		}
		responses.TextResponse(w, "Database connection is healthy")
	}
}
