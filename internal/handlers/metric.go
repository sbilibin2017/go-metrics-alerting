package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Интерфейс сервиса метрик
type MetricService interface {
	UpdateMetric(ctx context.Context, req *types.UpdateMetricValueRequest) error
	GetMetric(ctx context.Context, req *types.GetMetricValueRequest) (string, error)
	ListMetrics(ctx context.Context) []*types.MetricResponse
}

// Константа для HTML-шаблона
const metricsTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Metrics List</title>
</head>
<body>
	<h1>Metrics List</h1>
	<ul>
	{{range .}}
		<li>{{.Name}}: {{.Value}}</li>
	{{else}}
		<li>No metrics available</li>
	{{end}}
	</ul>
</body>
</html>`

func handleError(c *gin.Context, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		c.String(apiErr.Code, apiErr.Message)
	} else {
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func setHeaders(c *gin.Context, contentType string) {
	c.Header("Content-Type", contentType+"; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
}

// Обработчик обновления метрики
func UpdateMetricHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := service.UpdateMetric(c, &types.UpdateMetricValueRequest{
			Type:  types.MetricType(c.Param("type")),
			Name:  c.Param("name"),
			Value: c.Param("value"),
		})
		if err != nil {
			handleError(c, err)
			return
		}
		setHeaders(c, "text/plain")
		c.String(http.StatusOK, "Metric updated")
	}
}

// Обработчик получения метрики
func GetMetricHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricValueResp, err := service.GetMetric(c, &types.GetMetricValueRequest{
			Type: types.MetricType(c.Param("type")),
			Name: c.Param("name"),
		})
		if err != nil {
			handleError(c, err)
			return
		}
		setHeaders(c, "text/plain")
		c.String(http.StatusOK, metricValueResp)
	}
}

// Обработчик списка метрик
func ListMetricsHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := service.ListMetrics(c)
		tmpl, _ := template.New("metrics").Parse(metricsTemplate)
		setHeaders(c, "text/html")
		c.Writer.WriteHeader(http.StatusOK)
		tmpl.Execute(c.Writer, metrics)
	}
}
