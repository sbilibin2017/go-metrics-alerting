package main

import (
	"flag"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/pkg/logger"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
)

// Функция для запуска сервера
func main() {
	config := &configs.ServerConfig{}
	env.Parse(config)

	if config.Address == "" {
		config.Address = ":8080"
	}

	address := flag.String("a", "", "Address of the HTTP server (default: loaded from environment or :8080)")

	flag.Parse()

	if *address != "" {
		config.Address = *address
	}

	r := gin.Default()
	r.RedirectTrailingSlash = false

	storageEngine := &engines.StorageEngine{}

	keyEngine := &engines.KeyEngine{}

	metricRepository := &repositories.MetricRepository{
		StorageEngine: storageEngine,
		KeyEngine:     keyEngine,
	}

	metricService := &services.MetricService{MetricRepository: metricRepository}

	routers.RegisterMetricHandlers(r, metricService)

	if err := r.Run(config.Address); err != nil {
		logger.Logger.Fatalf("Failed to start server: %v", err)
	}
}
