package routers

import (
	"go-metrics-alerting/internal/responders"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Интерфейс сервиса метрик
type MetricService interface {
	UpdateMetric(req *types.UpdateMetricValueRequest)
	GetMetric(req *types.GetMetricValueRequest) (string, *types.APIErrorResponse)
	ListMetrics() []*types.MetricResponse
}

func RegisterRouter(r *gin.Engine, service MetricService) {
	r.POST("/update/:type/:name/:value", updateMetricHandler(service))
	r.GET("/value/:type/:name", getMetricHandler(service))
	r.GET("/", listMetricsHandler(service))
}

// Обработчик обновления метрики
func updateMetricHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorHandler := &responders.ErrorResponder{C: c}
		successHandler := &responders.SuccessResponder{C: c}

		// Создаем запрос из параметров
		req := &types.UpdateMetricValueRequest{
			Type:  types.MetricType(c.Param("type")),
			Name:  c.Param("name"),
			Value: c.Param("value"),
		}

		// Валидируем запрос с использованием метода Validate()
		if err := req.Validate(); err != nil {
			errorHandler.Respond(err)
			return
		}

		// Обновляем метрику
		service.UpdateMetric(req)
		successHandler.Respond("Metric updated", http.StatusOK)
	}
}

// Обработчик получения метрики
func getMetricHandler(service MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorHandler := &responders.ErrorResponder{C: c}
		successHandler := &responders.SuccessResponder{C: c}

		// Создаем запрос для получения метрики
		req := &types.GetMetricValueRequest{
			Type: types.MetricType(c.Param("type")),
			Name: c.Param("name"),
		}

		// Валидируем запрос с использованием метода Validate()
		if err := req.Validate(); err != nil {
			errorHandler.Respond(err)
			return
		}

		// Получаем метрику
		metricValueResp, err := service.GetMetric(req)
		if err != nil {
			errorHandler.Respond(err)
			return
		}

		// Успешный ответ
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
