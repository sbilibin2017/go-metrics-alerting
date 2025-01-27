package main

import (
	"fmt"
	"go-metrics-alerting/internal/configs" // Импортируем пакет конфигураций
	"go-metrics-alerting/internal/engines"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/routers/responders"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Импортируем godotenv для загрузки переменных окружения
)

func main() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Получаем конфигурацию сервера через флаги и переменные окружения
	serverConfig := configs.NewServerConfig()

	// Создаем движки для работы с метриками
	storageEngine := engines.NewMemoryStorageEngine()
	keyEngine := engines.NewKeyEngine()
	strategyEngines := map[types.MetricType]engines.StrategyUpdateEngineInterface{
		types.CounterType: &engines.CounterUpdateStrategyEngine{},
		types.GaugeType:   &engines.GaugeUpdateStrategyEngine{},
	}

	// Создаем сервис для работы с метриками
	metricService := services.NewMetricService(storageEngine, strategyEngines, keyEngine)

	// Пример использования:
	errorResponder := responders.NewErrorResponder()
	successResponder := responders.NewSuccessResponder()
	htmlResponder := responders.NewHTMLResponder()

	// Инициализируем роутер для метрик
	metricRouter := routers.NewMetricRouter(metricService, errorResponder, successResponder, htmlResponder)

	// Инициализация главного роутера
	router := initRouter()

	// Регистрируем маршруты
	metricRouter.RegisterMetricRoutes(router)

	// Запуск сервера
	router.Run(serverConfig.GetServerUrl())
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.RedirectFixedPath = false
	return router
}
