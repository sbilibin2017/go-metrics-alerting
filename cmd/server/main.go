package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"
	"go-metrics-alerting/internal/strategies"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логера
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Закрытие логера при завершении работы программы
	logger.Info("Logger initialized successfully")

	var config configs.ServerConfig
	parseFlagsAndEnv(&config)

	// Создание хранилища
	saver, getter, ranger := newStorage()

	// Создание стратегий обновления метрик
	updateGaugeStrategy, updateCounterStrategy := newUpdateStrategies(saver, getter)

	// Сервисы для обновления, получения метрик и получения всех метрик
	updateMetricService := services.NewUpdateMetricService(updateCounterStrategy, updateGaugeStrategy)
	getMetricService := services.NewGetMetricService(getter)
	getAllMetricsService := services.NewGetAllMetricsService(ranger)

	// Обработчики для маршрутов
	updateBodyHandler := handlers.UpdateMetricBodyHandler(updateMetricService)
	updatePathHandler := handlers.UpdateMetricPathHandler(updateMetricService)
	getBodyHandler := handlers.GetMetricBodyHandler(getMetricService)
	getPathHandler := handlers.GetMetricPathHandler(getMetricService)
	getAllHandler := handlers.GetAllMetricsHandler(getAllMetricsService)

	// Инициализация маршрутизатора
	r := chi.NewRouter()

	// Регистрируем обработчики
	routers.RegisterMetricsHandlers(
		r,
		logger,
		updateBodyHandler,
		updatePathHandler,
		getBodyHandler,
		getPathHandler,
		getAllHandler,
	)

	// Логируем запуск сервера
	logger.Info("Starting server on address", zap.String("address", config.Address))

	// Загружаем метрики при старте, если необходимо
	if config.Restore {
		err := loadMetricsFromFile(&config, updateMetricService)
		if err != nil {
			logger.Warn("Failed to load metrics from file", zap.Error(err))
		} else {
			logger.Info("Metrics restored from file")
		}
	}

	// Запускаем периодическое сохранение метрик
	go startPeriodicSaving(&config, getAllMetricsService, logger)

	// Запуск HTTP сервера
	err = http.ListenAndServe(config.Address, r)
	if err != nil {
		// Логирование ошибки, если сервер не запускается
		logger.Fatal("Server failed", zap.Error(err))
	}

	// Сохранение метрик при завершении работы
	defer func() {
		err := saveMetricsToFile(&config, getAllMetricsService, logger)
		if err != nil {
			logger.Error("Failed to save metrics before shutdown", zap.Error(err))
		} else {
			logger.Info("Metrics saved before shutdown")
		}
	}()
}

// Функция для инициализации логера
func newLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

// Функция для парсинга флагов и переменных окружения
func parseFlagsAndEnv(config *configs.ServerConfig) {
	// Конвертируем StoreInterval в секунды для флага
	storeIntervalInSeconds := int(config.StoreInterval.Seconds())

	// Флаги командной строки
	addressFlag := flag.String("address", "", "Address for the HTTP server")
	storeIntervalFlag := flag.Int("i", storeIntervalInSeconds, "Interval in seconds to store metrics to disk")
	fileStoragePathFlag := flag.String("f", config.FileStoragePath, "Path to the file where metrics will be saved")
	restoreFlag := flag.Bool("r", config.Restore, "Restore saved metrics from file on server start")

	// Парсим флаги
	flag.Parse()

	// Загружаем переменные окружения в структуру с помощью github.com/caarlos0/env
	err := env.Parse(config)
	if err != nil {
		fmt.Println("Error parsing environment variables:", err)
	}

	// Если флаг пустой, используем значение из переменной окружения
	if *addressFlag == "" && config.Address != "" {
		*addressFlag = config.Address
	}

	// Если флаг StoreInterval был передан, обновляем значение
	if *storeIntervalFlag != storeIntervalInSeconds {
		config.StoreInterval = time.Duration(*storeIntervalFlag) * time.Second
	}

	if *fileStoragePathFlag != config.FileStoragePath && config.FileStoragePath != "metrics_data.json" {
		*fileStoragePathFlag = config.FileStoragePath
	}

	if *restoreFlag != config.Restore {
		*restoreFlag = config.Restore
	}
}

// newStorage создает все необходимые структуры для работы с хранилищем данных.
func newStorage() (
	*storage.Saver[string, *domain.Metrics],
	*storage.Getter[string, *domain.Metrics],
	*storage.Ranger[string, *domain.Metrics],
) {
	s := storage.NewStorage[string, *domain.Metrics]()
	saver := storage.NewSaver(s)
	ranger := storage.NewRanger(s)
	getter := storage.NewGetter(s)
	return saver, getter, ranger
}

// Функция newUpdateStrategies для создания стратегий обновления метрик
func newUpdateStrategies(
	saver *storage.Saver[string, *domain.Metrics],
	getter *storage.Getter[string, *domain.Metrics],
) (*strategies.UpdateGaugeMetricStrategy, *strategies.UpdateCounterMetricStrategy) {
	updateGaugeStrategy := strategies.NewUpdateGaugeMetricStrategy(saver, getter)
	updateCounterStrategy := strategies.NewUpdateCounterMetricStrategy(saver, getter)
	return updateGaugeStrategy, updateCounterStrategy
}

func loadMetricsFromFile(config *configs.ServerConfig, service *services.UpdateMetricService) error {
	// Чтение метрик из файла
	data, err := os.ReadFile(config.FileStoragePath)
	if err != nil {
		return fmt.Errorf("failed to read metrics from file: %w", err)
	}

	// Десериализация метрик
	var metrics []*domain.Metrics
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	// Сохранение метрик в хранилище
	for _, metric := range metrics {
		service.UpdateMetric(metric)
	}

	return nil
}

func saveMetricsToFile(config *configs.ServerConfig, service *services.GetAllMetricsService, logger *zap.Logger) error {
	// Получаем все метрики
	metrics := service.GetAllMetrics()

	// Сериализуем метрики в формат JSON
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// Записываем данные в файл
	err = os.WriteFile(config.FileStoragePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metrics to file: %w", err)
	}

	return nil
}

func startPeriodicSaving(config *configs.ServerConfig, service *services.GetAllMetricsService, logger *zap.Logger) {
	// Интервал для сохранения метрик в файл
	ticker := time.NewTicker(config.StoreInterval)
	defer ticker.Stop()

	// Периодическое сохранение метрик
	for range ticker.C {
		err := saveMetricsToFile(config, service, logger)
		if err != nil {
			logger.Error("Failed to save metrics to file", zap.Error(err))
		} else {
			logger.Info("Metrics saved to file")
		}
	}
}
