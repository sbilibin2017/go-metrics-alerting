package routers

import (
	"context"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/types"

	"github.com/gin-gonic/gin"
)

// Интерфейс сервиса метрик
type MetricService interface {
	UpdateMetric(ctx context.Context, req *types.UpdateMetricValueRequest) error
	GetMetric(ctx context.Context, req *types.GetMetricValueRequest) (string, error)
	ListMetrics(ctx context.Context) []*types.MetricResponse
}

func RegisterMetricHandlers(r *gin.Engine, svc MetricService) {
	r.POST("/update/:type/:name/:value", handlers.UpdateMetricHandler(svc))
	r.GET("/value/:type/:name", handlers.GetMetricHandler(svc))
	r.GET("/", handlers.ListMetricsHandler(svc))

}
