package routers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/handlers"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsServiceInterface интерфейс для сервиса обновления метрик
type UpdateMetricsService interface {
	UpdateMetricValue(metric *domain.Metric) (*domain.Metric, error)
}

// GetMetricValueServiceInterface интерфейс для сервиса получения метрики по ID
type GetMetricValueService interface {
	GetMetricValue(id string, mType domain.MType) (*domain.Metric, error)
}

// GetAllMetricValuesServiceInterface интерфейс для сервиса получения всех метрик
type GetAllMetricValuesService interface {
	GetAllMetricValues() []*domain.Metric
}

// RegisterMetricRoutes регистрирует маршруты для метрик
func RegisterMetricRoutes(
	router *gin.Engine,
	updateMetricService UpdateMetricsService,
	getMetricValueService GetMetricValueService,
	getAllMetricValuesService GetAllMetricValuesService,
) {
	router.POST("/update/", handlers.UpdateMetricsBodyHandler(updateMetricService))
	// router.POST("/update/:mtype/:id/:value", handlers.UpdateMetricsPathHandler(updateMetricService))
	// router.POST("/value/", handlers.GetMetricValueBodyHandler(getMetricValueService))
	// router.GET("/value/:mtype/:id", handlers.GetMetricValuePathHandler(getMetricValueService))
	// router.GET("/", handlers.GetAllMetricValuesHandler(getAllMetricValuesService))
}
