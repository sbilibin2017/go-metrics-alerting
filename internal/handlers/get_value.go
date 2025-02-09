package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Интерфейс для получения значения метрики
type GetValueService interface {
	GetMetricValue(ctx context.Context, req *types.GetMetricValueRequest) (string, error)
}

// Регистрация обработчика для получения значения метрики
func RegisterGetMetricValueHandler(r *gin.Engine, svc GetValueService) {
	r.RedirectTrailingSlash = false

	r.GET("/value/:type/:name", func(c *gin.Context) {
		getMetricValueHandler(svc, c)
	})

	r.GET("/value/:type", func(c *gin.Context) {
		getMetricValueByTypeHandler(svc, c)
	})
}

// Обработчик получения значения метрики по типу и имени
func getMetricValueHandler(service GetValueService, c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	getRequest := &types.GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}

	// Передаем контекст в сервис
	metricValue, err := service.GetMetricValue(c, getRequest) // передаем контекст c
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	// Убеждаемся, что Content-Length устанавливается корректно
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.String(http.StatusOK, metricValue) // Gin сам установит Content-Length
}

// Обработчик получения значения метрики только по типу
func getMetricValueByTypeHandler(service GetValueService, c *gin.Context) {
	metricType := c.Param("type")

	getRequest := &types.GetMetricValueRequest{
		Type: metricType,
	}

	// Передаем контекст в сервис
	metricValue, err := service.GetMetricValue(c, getRequest) // передаем контекст c
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	c.String(http.StatusOK, metricValue) // Gin сам установит Content-Length
}
