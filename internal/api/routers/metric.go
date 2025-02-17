package routers

import (
	"go-metrics-alerting/internal/api/handlers"
	"go-metrics-alerting/internal/types"

	"github.com/gin-gonic/gin"
)

// UpdateMetricService интерфейс для сервиса обновления метрики.
type UpdateMetricService interface {
	UpdateMetric(req *types.MetricsRequest) (*types.MetricsResponse, error)
}

// NewMetricRouter инициализирует маршруты для метрик с префиксом "/" и возвращает роутер
func NewMetricRouter(svc UpdateMetricService) *gin.Engine {
	router := gin.Default()
	metricRouter := router.Group("/")
	metricRouter.POST("update", handlers.UpdateMetricHandler(svc))
	return router
}
