package routers

import (
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/types"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsServiceInterface интерфейс для сервиса обновления метрик
type UpdateMetricsService interface {
	UpdateMetricValue(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, *types.APIErrorResponse)
	ParseMetricValues(mtype, valueStr string) (*float64, *int64, *types.APIErrorResponse)
}

// GetMetricValueServiceInterface интерфейс для сервиса получения метрики по ID
type GetMetricValueService interface {
	GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, *types.APIErrorResponse)
}

// GetAllMetricValuesServiceInterface интерфейс для сервиса получения всех метрик
type GetAllMetricValuesService interface {
	GetAllMetricValues() ([]*types.GetMetricValueResponse, *types.APIErrorResponse)
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
