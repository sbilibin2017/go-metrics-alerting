package server

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/middlewares"
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Run() {
	config := &configs.ServerConfig{}

	godotenv.Load()
	env.Parse(config)

	if config.Address == "" {
		config.Address = ":8080"
	}

	address := flag.String("a", "", "Address of the HTTP server (default: loaded from environment or :8080)")

	flag.Parse()

	if *address != "" {
		config.Address = *address
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := gin.Default()
	r.RedirectTrailingSlash = false

	r.Use(middlewares.LoggerMiddleware(logger))

	memStorage := storage.NewMemStorage()

	metricRepo := repositories.NewMetricRepository(memStorage)

	metricService := services.NewMetricService(metricRepo)

	routers.RegisterRouter(r, metricService)

	r.Run(config.Address)
}
