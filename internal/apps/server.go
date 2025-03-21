package apps

import (
	"context"
	"database/sql"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultServerAddress   = ":8080"
	DefaultDatabaseDSN     = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
	DefaultStoreInterval   = ""
	DefaultFileStoragePath = ""
	DefaultRestore         = ""

	EnvServerAddress   = "ADDRESS"
	EnvDatabaseDSN     = "DATABASE_DSN"
	EnvStoreInterval   = "STORE_INTERVAL"
	EnvFileStoragePath = "FILE_STORAGE_PATH"
	EnvRestore         = "RESTORE"

	FlagServerAddress   = "a"
	FlagDatabaseDSN     = "d"
	FlagStoreInterval   = "s"
	FlagFileStoragePath = "f"
	FlagRestore         = "r"

	DescriptionServerAddress   = "Server address"
	DescriptionDatabaseDSN     = "Database DSN"
	DescriptionStoreInterval   = "Interval in seconds for data store"
	DescriptionFileStoragePath = "Path to file storage"
	DescriptionRestore         = "Restore backup (true/false)"
)

// NewServerCommand initializes the Cobra command for the server configuration.
func NewServerCommand() *cobra.Command {
	var config configs.ServerConfig

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Initialize server configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Retrieve configuration values
			config.Address = viper.GetString(FlagServerAddress)
			config.DatabaseDSN = viper.GetString(FlagDatabaseDSN)
			config.StoreInterval = viper.GetString(FlagStoreInterval)
			config.FileStoragePath = viper.GetString(FlagFileStoragePath)
			config.Restore = viper.GetString(FlagRestore)

			// Set defaults for missing config values
			if config.Address == "" {
				config.Address = DefaultServerAddress
			}
			if config.DatabaseDSN == "" {
				config.DatabaseDSN = DefaultDatabaseDSN
			}
			if config.StoreInterval == "" {
				config.StoreInterval = DefaultStoreInterval
			}
			if config.FileStoragePath == "" {
				config.FileStoragePath = DefaultFileStoragePath
			}
			if config.Restore == "" {
				config.Restore = DefaultRestore
			}

			// Set up signal context for graceful shutdown
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			return runServerApp(ctx, &config)
		},
	}

	// Define flags
	cmd.Flags().String(FlagServerAddress, DefaultServerAddress, DescriptionServerAddress)
	cmd.Flags().String(FlagDatabaseDSN, DefaultDatabaseDSN, DescriptionDatabaseDSN)
	cmd.Flags().String(FlagStoreInterval, DefaultStoreInterval, DescriptionStoreInterval)
	cmd.Flags().String(FlagFileStoragePath, DefaultFileStoragePath, DescriptionFileStoragePath)
	cmd.Flags().String(FlagRestore, DefaultRestore, DescriptionRestore)

	// Bind flags to Viper
	viper.BindPFlag(FlagServerAddress, cmd.Flags().Lookup(FlagServerAddress))
	viper.BindPFlag(FlagDatabaseDSN, cmd.Flags().Lookup(FlagDatabaseDSN))
	viper.BindPFlag(FlagStoreInterval, cmd.Flags().Lookup(FlagStoreInterval))
	viper.BindPFlag(FlagFileStoragePath, cmd.Flags().Lookup(FlagFileStoragePath))
	viper.BindPFlag(FlagRestore, cmd.Flags().Lookup(FlagRestore))

	// Set up Viper to read environment variables automatically
	viper.AutomaticEnv()

	// Bind environment variables to Viper
	viper.BindEnv(FlagServerAddress, EnvServerAddress)
	viper.BindEnv(FlagDatabaseDSN, EnvDatabaseDSN)
	viper.BindEnv(FlagStoreInterval, EnvStoreInterval)
	viper.BindEnv(FlagFileStoragePath, EnvFileStoragePath)
	viper.BindEnv(FlagRestore, EnvRestore)

	return cmd
}

// runServerApp creates and initializes the server with the provided configuration.
func runServerApp(ctx context.Context, config *configs.ServerConfig) error {
	var file *os.File
	var db *sql.DB
	var err error

	// Assuming `config.FileStoragePath` contains the full path to the file
	if config.FileStoragePath != "" {
		// Create all necessary directories leading to the file
		dir := filepath.Dir(config.FileStoragePath)
		err := os.MkdirAll(dir, 0755) // 0755 is a common permission setting for directories
		if err != nil {
			return err
		}

		// Open the file for reading and writing (if the file doesn't exist, it will be created)
		file, err = os.OpenFile(config.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	// 2. Connect to the database if DatabaseDSN is provided in the config
	if config.DatabaseDSN != "" {
		// Open a connection to the database using pgx
		db, err = sql.Open("pgx", config.DatabaseDSN)
		if err != nil {
			return err
		}
		// Check if the database is reachable
		if err := db.PingContext(ctx); err != nil {
			return err
		}
		defer db.Close() // Ensure the DB connection is closed when the function exits
	}

	// 3. Repositories
	metricRepo := repositories.NewMetricRepository(config, file, db)

	metricService := services.NewMetricService(metricRepo.GetMainRepository(config))

	// Create a new router
	r := chi.NewRouter()

	// 8. Set up the /ping route to check DB health
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := db.PingContext(ctx)
		if err != nil {
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Database connection successful")
	})

	// 9. Set up the /metrics route and other routes for the metric handler
	metricHandler := handlers.NewMetricHandler(metricService)
	metricRouter := routers.NewMetricRouter(config, metricHandler)
	r.Mount("/", metricRouter) // Mount the metric router

	// 10. Initialize the HTTP server
	server := &http.Server{
		Addr:    config.Address,
		Handler: r, // Attach the router to the server
	}

	// Run server and workers
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cancelCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			<-cancelCtx.Done() // En
		}
	}()

	// Start workers directly without using WorkerRegistry
	go func() {
		if err := loadMetricsFromFile(ctx, config, metricRepo); err != nil {
			cancelCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			<-cancelCtx.Done() // En
		}
	}()

	go func() {
		if err := dumpMetricsToFile(ctx, config, metricRepo); err != nil {
			cancelCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			<-cancelCtx.Done() // En
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error during server shutdown: %w", err)
	}

	return nil
}

// LoadMetricsFromFile loads metrics from a secondary repository (e.g., file storage)
// and saves them to the main repository if restoration is enabled in the configuration.
func loadMetricsFromFile(
	ctx context.Context,
	config *configs.ServerConfig,
	repo *repositories.MetricRepository,
) error {
	// If no file storage path is specified, do nothing
	if config.FileStoragePath == "" {
		return nil
	}

	// If restore is disabled in the config, do nothing
	if config.Restore == "" || config.Restore == "false" {
		return nil
	}

	// Retrieve metrics from the secondary repository (e.g., file storage)
	metrics, err := repo.FileRepo.ListMetrics(ctx)
	if err != nil {
		return err
	}

	// If there are metrics, save them to the main repository
	if len(metrics) != 0 {
		err = repo.GetMainRepository(config).SaveMetrics(ctx, metrics)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpMetricsToFile(
	ctx context.Context,
	config *configs.ServerConfig,
	repo *repositories.MetricRepository,
) error {
	// Define the work function for dumping metrics
	work := func() error {
		metrics, err := repo.GetMainRepository(config).ListMetrics(ctx)
		if err != nil {
			return err
		}

		if len(metrics) != 0 {
			// Save the metrics to the secondary repository (e.g., file storage)
			err = repo.FileRepo.SaveMetrics(ctx, metrics)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// Get the store interval from the configuration
	storeIntervalStr := config.StoreInterval
	if storeIntervalStr == "0" || storeIntervalStr == "" {
		// If no interval is configured, execute the work immediately and only once
		err := work()
		if err != nil {
			return err
		}
		return nil
	} else {
		// Parse the store interval to an integer (in seconds)
		storeInterval, err := strconv.Atoi(storeIntervalStr)
		if err != nil || storeInterval <= 0 {
			return fmt.Errorf("invalid store interval: %v", storeIntervalStr)
		}

		// Set up a ticker to execute the work at the specified interval
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		defer ticker.Stop()

		// Continuously run the work function at the defined intervals
		for {
			select {
			case <-ctx.Done():
				// If the context is done, perform the work one last time and return
				err := work()
				if err != nil {
					return err
				}
				return ctx.Err() // Return the context error (if canceled or timed out)
			case <-ticker.C:
				// Run the work function at every tick
				err := work()
				if err != nil {
					return err
				}
			}
		}
	}
}
