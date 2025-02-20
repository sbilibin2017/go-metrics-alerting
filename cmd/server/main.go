package main

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/keyprocessor"
	"go-metrics-alerting/internal/routers"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/storage"
	"go-metrics-alerting/internal/strategies"
)

func main() {
	// Инициализация хранилища
	strg := storage.NewStorage()

	// Инициализация Saver, Getter и Ranger через их конструкторы
	saver := storage.NewSaver(strg)   // Создание нового Saver
	getter := storage.NewGetter(strg) // Создание нового Getter
	ranger := storage.NewRanger(strg) // Создание нового Ranger

	keyEncoder := keyprocessor.NewKeyEncoder()
	keyDecoder := keyprocessor.NewKeyDecoder()

	// Инициализация стратегий с использованием конструктора
	updateGaugeStrategy := strategies.NewUpdateGaugeStrategy(saver, keyEncoder)
	updateCounterStrategy := strategies.NewUpdateCounterStrategy(saver, getter, keyEncoder)
	updateStrategies := map[domain.MType]services.UpdateMetricStrategy{
		domain.Gauge:   updateGaugeStrategy,
		domain.Counter: updateCounterStrategy,
	}

	updateMetricsService := services.NewUpdateMetricsService(updateStrategies)
	getMetricValueService := services.NewGetMetricValueService(getter, keyEncoder)
	getAllMetricValuesService := services.NewGetAllMetricValuesService(ranger, keyDecoder)

	h1 := handlers.UpdateMetricsBodyHandler(updateMetricsService)
	h2 := handlers.UpdateMetricsPathHandler(updateMetricsService)
	h3 := handlers.GetMetricValueBodyHandler(getMetricValueService)
	h4 := handlers.GetMetricValuePathHandler(getMetricValueService)
	h5 := handlers.GetAllMetricValuesHandler(getAllMetricValuesService)

	routers.RegisterMetricRoutes(r, h1, h2, h3, h4, h5)

	r.Run()
}
