package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateValueService interface {
	UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error
}

// Регистрация обработчиков обновления метрик
func RegisterUpdateValueHandler(r *gin.Engine, svc UpdateValueService) {
	r.RedirectTrailingSlash = false

	r.POST("/update/:type/:name/:value", func(c *gin.Context) {
		updateValueHandler(svc, c)
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Route not found")
	})
}

// Обработчик обновления значения метрики
func updateValueHandler(service UpdateValueService, c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	updateRequest := &types.UpdateMetricValueRequest{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}

	// Обновляем значение метрики с валидацией в сервисе
	if err := service.UpdateMetricValue(c, updateRequest); err != nil { // Передаем контекст c
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
			return
		}
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Формируем ответ
	response := "Metric updated"
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.String(http.StatusOK, response) // Gin сам установит Content-Length
}
