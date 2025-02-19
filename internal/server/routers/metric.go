package routers

import (
	"go-metrics-alerting/internal/server/handlers"
	"go-metrics-alerting/internal/server/types"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsService определяет контракт для обновления метрик.
type UpdateMetricsService interface {
	Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error)
}

// GetMetricValueService определяет контракт для получения значения метрики.
type GetMetricValueService interface {
	GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error)
}

// GetAllMetricValuesService интерфейс для сервиса работы с метриками.
type GetAllMetricValuesService interface {
	GetAllMetricValues() []*types.GetMetricValueResponse
}

// RegisterMetricRoutes регистрирует маршруты для метрик
func RegisterMetricRoutes(
	router *gin.Engine,
	updateMetricsService UpdateMetricsService,
	getMetricValueService GetMetricValueService,
	getAllMetricValuesService GetAllMetricValuesService,
) {
	router.POST("/update/", handlers.UpdateMetricsBodyHandler(updateMetricsService))
	router.POST("/update/:mtype/:id/:value", handlers.UpdateMetricsPathHandler(updateMetricsService))
	router.POST("/value/", handlers.GetMetricValueBodyHandler(getMetricValueService))
	router.GET("/value/:mtype/:id", handlers.GetMetricValuePathHandler(getMetricValueService))
	router.GET("/", handlers.GetAllMetricValuesHandler(getAllMetricValuesService))
}
