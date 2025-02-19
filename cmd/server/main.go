package main

import (
	"flag"
	"go-metrics-alerting/internal/server/configs"
	"go-metrics-alerting/internal/server/logger"
	"go-metrics-alerting/internal/server/middlewares"
	"go-metrics-alerting/internal/server/routers"
	"go-metrics-alerting/internal/server/services"
	"go-metrics-alerting/internal/server/storage"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	godotenv.Load()

	// Создание конфигурации
	var config configs.ServerConfig
	env.Parse(&config)

	// Проверка наличия флага для адреса
	addressFlag := flag.String("address", "", "Address for HTTP server")
	flag.Parse()

	// Если адрес из .env не установлен, используем флаг
	if config.Address == "" {
		config.Address = *addressFlag
	}

	// Создание нового gin сервера
	r := gin.New()
	r.RedirectTrailingSlash = false

	// Мидлваре
	r.Use(middlewares.JSONContentTypeMiddleware())
	r.Use(middlewares.LoggerMiddleware(logger.Logger))

	// Инициализация хранилища и сервисов
	s := storage.NewStorage()
	saver := storage.NewSaver(s)
	getter := storage.NewGetter(s)
	ranger := storage.NewRanger(s)

	updateMetricsService := services.NewUpdateMetricsService(saver, getter)
	getMetricValueService := services.NewGetMetricValueService(getter)
	getAllMetricValuesService := services.NewGetAllMetricValuesService(ranger)

	// Регистрация маршрутов
	routers.RegisterMetricRoutes(
		r,
		updateMetricsService,
		getMetricValueService,
		getAllMetricValuesService,
	)

	// Запуск сервера
	r.Run(config.Address)
}
