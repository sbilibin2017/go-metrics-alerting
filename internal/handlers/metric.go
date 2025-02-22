package handlers

import (
	"fmt"
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/templates"
	"go-metrics-alerting/internal/types"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateMetricService фасад для двух сервисов.
type UpdateMetricService interface {
	UpdateMetric(metric *domain.Metrics) *domain.Metrics
}

// GetMetricService для получения одной метрики по её ID.
type GetMetricService interface {
	GetMetric(id string, mtype domain.MType) *domain.Metrics
}

// GetAllMetricsService для получения всех метрик с использованием Ranger.
type GetAllMetricsService interface {
	GetAllMetrics() []*domain.Metrics
}

// Обработчик для обновления метрики через тело запроса
func UpdateMetricsBodyHandler(svc UpdateMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric types.UpdateMetricBodyRequest
		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if err := metric.Validate(); err != nil {
			c.JSON(err.Status, gin.H{"error": err.Message})
			return
		}

		domainMetric := metric.ToDomain()
		updatedMetric := svc.UpdateMetric(domainMetric)
		c.JSON(http.StatusOK, updatedMetric)
	}
}

// Обработчик для обновления метрики через параметры пути запроса
func UpdateMetricsPathHandler(svc UpdateMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric types.UpdateMetricPathRequest
		metric.ID = c.Param("id")
		metric.MType = c.Param("mtype")
		metric.Value = c.Param("value")

		fmt.Println("Received parameters:", metric)
		if err := metric.Validate(); err != nil {
			fmt.Println("Validation error:", err.Message)
			c.String(err.Status, err.Message)
			return
		}

		domainMetric := metric.ToDomain()
		updatedMetric := svc.UpdateMetric(domainMetric)
		if updatedMetric == nil {
			fmt.Println("Failed to update metric")
			c.String(http.StatusBadRequest, "Metric is not updated")
			return
		}
		c.String(http.StatusOK, "Metric updated successfully")
	}
}

// Обработчик для получения метрики через тело запроса
func GetMetricValueBodyHandler(svc GetMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.GetMetricRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if err := request.Validate(); err != nil {
			c.JSON(err.Status, gin.H{"error": err.Message})
			return
		}

		metric := svc.GetMetric(request.ID, domain.MType(request.MType))
		if metric == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("%v", metric))
	}
}

// Обработчик для получения метрики через параметры пути запроса
func GetMetricValuePathHandler(svc GetMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		mType := domain.MType(c.Param("mType"))

		metric := svc.GetMetric(id, mType)
		if metric == nil {
			c.String(http.StatusNotFound, "Metric not found")
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("%f", *metric.Value))
	}
}

func GetAllMetricValuesHandler(svc GetAllMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tmpl, _ := template.New("metrics").Parse(templates.MetricsTemplate)
		metrics := svc.GetAllMetrics()
		tmpl.Execute(c.Writer, metrics)
	}
}
