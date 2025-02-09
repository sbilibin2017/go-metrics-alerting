package main

import (
	"flag"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/repositories"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/pkg/logger"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
)

// Функция для запуска сервера
func main() {
	config := &configs.ServerConfig{}
	if err := env.Parse(config); err != nil {

	}

	// Если переменная окружения не задана, устанавливаем значение по умолчанию
	if config.Address == "" {
		config.Address = ":8080"
	}

	// Определяем флаг для адреса сервера, с дефолтным значением
	address := flag.String("a", "", "Address of the HTTP server (default: loaded from environment or :8080)")

	// Парсим флаги
	flag.Parse()

	// Если флаг адреса был задан, переопределяем значение
	if *address != "" {
		config.Address = *address
	}

	// Создаем новый экземпляр Gin
	r := gin.Default()
	r.RedirectTrailingSlash = false

	// Создаем хранилище данных (в данном случае это память, но может быть база данных)
	storageEngine := &engines.StorageEngine{}

	// Создаем обработчик ключей для хранилища
	keyEngine := &engines.KeyEngine{}

	// Создаем репозиторий для метрик
	metricRepository := &repositories.MetricRepository{
		StorageEngine: storageEngine,
		KeyEngine:     keyEngine,
	}

	// Создаем сервисы для работы с метриками
	updateMetricService := &services.UpdateMetricValueService{MetricRepository: metricRepository}
	getMetricService := &services.GetMetricValueService{MetricRepository: metricRepository}
	getAllMetricService := &services.GetAllMetricValuesService{MetricRepository: metricRepository}

	// Регистрируем обработчики для маршрутов
	handlers.RegisterUpdateValueHandler(r, updateMetricService)
	handlers.RegisterGetMetricValueHandler(r, getMetricService)
	handlers.RegisterGetAllMetricValuesHandler(r, getAllMetricService)

	// Логируем информацию о запуске сервера
	logger.Logger.Infof("Server is running on %s...", config.Address)

	// Запуск сервера на указанном порту
	if err := r.Run(config.Address); err != nil {
		// Логируем ошибку, если сервер не может запуститься
		logger.Logger.Fatalf("Failed to start server: %v", err)
	}
}
