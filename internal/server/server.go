package server

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/logger" // Import the logger package
	"go-metrics-alerting/internal/middlewares"
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() {
	// Load configuration
	config := &configs.ServerConfig{}
	godotenv.Load()
	env.Parse(config)

	// Set default server address if not set
	if config.Address == "" {
		config.Address = ":8080"
	}

	address := flag.String("a", "", "Address of the HTTP server (default: loaded from environment or :8080)")
	flag.Parse()

	if *address != "" {
		config.Address = *address
	}

	// Initialize logger configuration
	logConfig := &configs.LoggerConfig{
		LogLevel: configs.INFO, // Or any other log level you'd like to configure
	}

	// Create the logger instance based on the configuration
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Create the Gin engine
	r := gin.Default()
	r.RedirectTrailingSlash = false

	// Pass the logger to the LoggerMiddleware
	r.Use(middlewares.LoggerMiddleware(log)) // Pass the logger instance

	// Initialize in-memory storage
	memStorage := storage.NewMemStorage()

	// Initialize repository and service
	metricRepo := repositories.NewMetricRepository(memStorage)
	metricService := services.NewMetricService(metricRepo)

	// Register the router with the metric service
	routers.RegisterRouter(r, metricService)

	// Start the server
	r.Run(config.Address)
}
