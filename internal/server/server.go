package server

import (
	"flag"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/middlewares"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// NewServer создает новый сервер с роутером и middleware
func RunServer() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		logger.Logger.Warn("Error loading .env file")
	}

	// Чтение конфигурации с использованием пакета github.com/caarlos0/env
	var config configs.ServerConfig
	if err := env.Parse(&config); err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to parse environment variables: %v", err))
	}

	// Обработка флага командной строки для адреса
	addressFlag := flag.String("address", "", "Address for HTTP server")
	flag.Parse()

	// Приоритет значений:
	// 1. Переменная окружения (если указана)
	// 2. Флаг командной строки (если указан)
	// 3. Значение по умолчанию
	if *addressFlag != "" {
		config.Address = *addressFlag
	}

	// Логирование адреса
	logger.Logger.Info(fmt.Sprintf("Server will run on address: %s", config.Address))

	// Инициализация сервера
	r := gin.New()
	r.RedirectTrailingSlash = false

	// Подключение миддлваров
	r.Use(middlewares.JSONContentTypeMiddleware())
	r.Use(middlewares.LoggerMiddleware(logger.Logger))

	// Инициализация хранилищ для гейджов и счётчиков
	s := storage.NewStorage()

	// Инициализация репозиториев
	saver := storage.NewSaver(s)
	getter := storage.NewGetter(s)
	ranger := storage.NewRanger(s)

	// Инициализация сервисов
	updateMetricsService := services.NewUpdateMetricsService(
		saver, getter,
	)

	getMetricValueService := services.NewGetMetricValueService(
		getter,
	)

	getAllMetricValuesService := services.NewGetAllMetricValuesService(
		ranger,
	)

	// Регистрация маршрутов для получения метрик
	routers.RegisterMetricRoutes(
		r,
		updateMetricsService,
		getMetricValueService,
		getAllMetricValuesService,
	)

	// Запуск сервера на нужном адресе
	r.Run(config.Address)
}
