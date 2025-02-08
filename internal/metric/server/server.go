package server

import (
	"flag"
	"go-metrics-alerting/internal/metric/configs"
	"go-metrics-alerting/internal/metric/handlers"
	"go-metrics-alerting/internal/metric/services"
	"go-metrics-alerting/internal/storage"
	"go-metrics-alerting/internal/storage/repositories"
	"go-metrics-alerting/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Функция для запуска сервера
func Run() {
	// Определяем флаг для адреса сервера, с дефолтным значением
	address := flag.String("a", "", "Address of the HTTP server (default: loaded from environment or :8080)")

	// Парсим флаги
	flag.Parse()

	// Загружаем конфигурацию сервера, которая будет учитывать переменные окружения
	config, err := configs.LoadServerConfig()
	if err != nil {
		logger.Logger.Fatalf("Failed to load server config: %v", err)
	}

	// Если флаг адреса был задан, переопределяем значение
	if *address != "" {
		config.Address = *address
	}

	// Создаем новый экземпляр Gin
	r := gin.Default()
	r.RedirectTrailingSlash = false

	// Создаем хранилище данных (в данном случае это память, но может быть база данных)
	memStorage := storage.NewMemStorage()

	// Создаем обработчик ключей для хранилища
	keyProcessor := storage.NewKeyProcessor()

	// Создаем репозиторий для метрик
	metricRepository := repositories.NewMetricRepository(memStorage, keyProcessor)

	// Создаем сервисы для работы с метриками
	updateMetricService := services.NewUpdateMetricValueService(metricRepository)
	getMetricService := services.NewGetMetricValueService(metricRepository)
	getAllMetricService := services.NewGetAllMetricsService(metricRepository)

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
