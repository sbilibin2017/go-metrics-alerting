package main

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/keyencoder"
	"go-metrics-alerting/internal/middlewares"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"
	"go-metrics-alerting/internal/strategies"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Инициализация конфигурации
	var config configs.ServerConfig
	if err := env.Parse(&config); err != nil {
		panic("Failed to parse environment variables: " + err.Error())
	}

	// Парсинг флагов командной строки
	addressFlag := flag.String("address", "", "Address for the HTTP server")
	flag.Parse()

	// Переопределение адреса с флага, если указан
	if *addressFlag != "" {
		config.Address = *addressFlag
	}

	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Логируем успешную инициализацию логгера
	logger.Info("Logger initialized")

	// Инициализация базы данных и других компонентов
	db := storage.NewStorage()
	saver := storage.NewSaver(db)
	getter := storage.NewGetter(db)
	ranger := storage.NewRanger(db)
	encoder := keyencoder.NewKeyEncoder()

	// Логируем создание базы данных и компонентов
	logger.Info("Database and components initialized")

	// Инициализация стратегий обновления метрик
	counterUpdateStrategy := strategies.NewUpdateCounterMetricStrategy(saver, getter, encoder, logger)
	gaugeUpdateStrategy := strategies.NewUpdateGaugeMetricStrategy(saver, getter, encoder, logger)

	// Логируем успешную инициализацию стратегий
	logger.Info("Metric update strategies initialized")

	// Инициализация сервисов
	updateMetricService := services.NewUpdateMetricService(counterUpdateStrategy, gaugeUpdateStrategy, logger)
	getMetricService := services.NewGetMetricService(getter, encoder, logger)
	getAllMetricService := services.NewGetAllMetricsService(ranger, logger)

	// Логируем успешную инициализацию сервисов
	logger.Info("Services initialized")

	// Инициализация маршрутов и middleware
	r := gin.New()
	r.RedirectTrailingSlash = false

	// Использование логирования и других middleware
	r.Use(middlewares.LoggingMiddleware(logger))
	r.Use(middlewares.SetHeadersMiddleware())

	// Регистрируем маршруты
	routers.RegisterMetricRoutes(
		r,
		handlers.UpdateMetricsBodyHandler(updateMetricService),
		handlers.UpdateMetricsPathHandler(updateMetricService),
		handlers.GetMetricValueBodyHandler(getMetricService),
		handlers.GetMetricValuePathHandler(getMetricService),
		handlers.GetAllMetricValuesHandler(getAllMetricService),
	)

	// Логируем успешную регистрацию маршрутов
	logger.Info("Routes registered")

	// Запуск сервера с заданным адресом
	logger.Info("Server started", zap.String("Address", config.Address))
	if err := r.Run(config.Address); err != nil {
		// Логируем ошибку запуска сервера
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
