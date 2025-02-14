package routers

import (
	"go-metrics-alerting/internal/responders"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/internal/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Интерфейс сервиса метрик
type MetricService interface {
	UpdateMetric(req *types.UpdateMetricValueRequest)
	GetMetric(req *types.GetMetricValueRequest) (string, *types.APIErrorResponse)
	ListMetrics() []*types.MetricResponse
}

func NewMetricRouter(service MetricService) *gin.Engine {
	r := gin.Default()
	r.POST("/update/:type/:name/:value", updateMetricHandler(service))
	r.GET("/value/:type/:name", getMetricHandler(service))
	r.GET("/", listMetricsHandler(service))
	return r
}

// Обработчик обновления метрики
func updateMetricHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorHandler := &responders.ErrorResponder{C: c}
		successHandler := &responders.SuccessResponder{C: c}

		req := &types.UpdateMetricValueRequest{
			Type:  types.MetricType(c.Param("type")),
			Name:  c.Param("name"),
			Value: c.Param("value"),
		}

		if err := validators.ValidateUpdateMetricRequest(req); err != nil {
			errorHandler.Respond(err)
			return
		}
		if err := validators.ValidateMetricValue(req.Value, req.Type); err != nil {
			errorHandler.Respond(err)
			return
		}

		service.UpdateMetric(req)
		successHandler.Respond("Metric updated", http.StatusOK)
	}
}

// Обработчик получения метрики
func getMetricHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorHandler := &responders.ErrorResponder{C: c}
		successHandler := &responders.SuccessResponder{C: c}
		metricValueResp, err := service.GetMetric(&types.GetMetricValueRequest{
			Type: types.MetricType(c.Param("type")),
			Name: c.Param("name"),
		})
		if err != nil {
			errorHandler.Respond(err)
			return
		}
		successHandler.Respond(metricValueResp, http.StatusOK)
	}
}

// Обработчик списка метрик
func listMetricsHandler(service MetricService) gin.HandlerFunc {
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
	return func(c *gin.Context) {
		htmlHandler := &responders.HTMLHandler{C: c}
		metrics := service.ListMetrics()
		htmlHandler.Respond(metricsTemplate, metrics)
	}
}
