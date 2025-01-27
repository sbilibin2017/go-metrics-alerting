package routers

import (
	"go-metrics-alerting/internal/routers/responders"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MetricUpdateRoute = "/update/:type/:name/:value"
	MetricValueRoute  = "/value/:type/:name"
	AllMetricsRoute   = "/"
)

// MetricRouter это структура, которая инкапсулирует зависимости для работы с роутами
type MetricRouter struct {
	service          services.MetricServiceInterface
	errorResponder   responders.ResponderInterface
	successResponder responders.ResponderInterface
	htmlResponder    responders.ResponderInterface
}

// NewMetricRouter создает новый экземпляр MetricRouter
// Все респондеры теперь передаются через аргументы
func NewMetricRouter(
	service services.MetricServiceInterface,
	errorResponder responders.ResponderInterface,
	successResponder responders.ResponderInterface,
	htmlResponder responders.ResponderInterface,
) *MetricRouter {
	return &MetricRouter{
		service:          service,
		errorResponder:   errorResponder,
		successResponder: successResponder,
		htmlResponder:    htmlResponder,
	}
}

// RegisterMetricRoutes регистрирует маршруты для метрик
func (r *MetricRouter) RegisterMetricRoutes(router *gin.Engine) {
	router.POST(MetricUpdateRoute, r.updateMetric)
	router.GET(MetricValueRoute, r.getMetric)
	router.GET(AllMetricsRoute, r.getAllMetrics)
}

// updateMetric обрабатывает обновление метрики
func (r *MetricRouter) updateMetric(c *gin.Context) {
	metricType, metricName, metricValue := extractParams(c)
	err := r.service.UpdateMetric(types.UpdateMetricRequest{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	})

	if err != nil {
		r.errorResponder.Respond(c, err.Status(), err.Error())
		return
	}

	r.successResponder.Respond(c, http.StatusOK, "Metric updated successfully")
}

// getMetric обрабатывает запрос на получение метрики
func (r *MetricRouter) getMetric(c *gin.Context) {
	metricType, metricName, _ := extractParams(c)
	value, err := r.service.GetMetric(types.GetMetricRequest{
		Type: metricType,
		Name: metricName,
	})

	if err != nil {
		r.errorResponder.Respond(c, err.Status(), err.Error())
		return
	}

	r.successResponder.Respond(c, http.StatusOK, value)
}

// getAllMetrics обрабатывает запрос на получение всех метрик
func (r *MetricRouter) getAllMetrics(c *gin.Context) {
	metrics := r.service.GetAllMetrics()
	htmlResponse := "<html><body><h1>Metrics</h1><ul>"
	for _, metric := range metrics {
		htmlResponse += "<li>" + metric.Type + " " + metric.Name + ": " + metric.Value + "</li>"
	}
	htmlResponse += "</ul></body></html>"
	r.htmlResponder.Respond(c, http.StatusOK, htmlResponse)
}

// extractParams извлекает параметры из запроса
func extractParams(c *gin.Context) (string, string, string) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value") // Это значение может быть пустым для некоторых запросов
	return metricType, metricName, metricValue
}
