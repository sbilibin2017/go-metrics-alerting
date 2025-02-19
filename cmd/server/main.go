package main

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/middlewares"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"
	"log"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
)

func main() {
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

	// Мидлвар
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

	log.Printf("Starting server on %s", config.Address)

	// Запуск сервера
	r.Run(config.Address)

}
