package handlers

import (
	"go-metrics-alerting/internal/types"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

// UpdateMetricService интерфейс для сервиса обновления метрики.
type UpdateMetricService interface {
	UpdateMetric(req *types.MetricsRequest) (*types.MetricsResponse, *types.APIStatusResponse)
}

func UpdateMetricHandler(svc UpdateMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Определяем метрику как указатель
		req := &types.MetricsRequest{}

		// Пробуем привязать JSON в тело запроса к структуре metric
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Логика обновления метрики через сервис
		metric, err := svc.UpdateMetric(req)
		if err != nil {
			// Если ошибка при обновлении метрики, используем err.Status и err.Message
			c.JSON(err.Status, gin.H{
				"error": err.Message,
			})
			return
		}

		// Возвращаем обновленную метрику
		c.JSON(http.StatusOK, metric)
	}
}

// GetMetricService интерфейс для сервиса получения метрики.
type GetMetricService interface {
	GetMetric(req *types.MetricValueRequest) (*string, *types.APIStatusResponse)
}

func GetMetricHandler(svc GetMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Определяем метрику как указатель
		req := &types.MetricValueRequest{}

		// Пробуем привязать JSON в тело запроса к структуре MetricValueRequest
		if err := c.ShouldBindJSON(req); err != nil {
			// Если ошибка при привязке, возвращаем ошибку с кодом 400
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Получаем метрику через сервис
		value, err := svc.GetMetric(req)
		if err != nil {
			// Если произошла ошибка (например, метрика не найдена), возвращаем ошибку
			c.JSON(err.Status, gin.H{
				"error": err.Message,
			})
			return
		}

		// Если метрика найдена, возвращаем её значение
		c.JSON(http.StatusOK, gin.H{
			"id":    req.ID,
			"mtype": req.MType,
			"value": *value,
		})
	}
}

// MetricService интерфейс для сервиса работы с метриками.
type ListMetricsService interface {
	// ListMetrics возвращает список всех метрик.
	ListMetrics() []*types.MetricValueResponse
}

// ListMetricsHandler обрабатывает запрос и возвращает HTML-страницу со списком метрик
func ListMetricsHandler(service ListMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем список метрик из сервиса
		metrics := service.ListMetrics()

		// Создаем и парсим HTML-шаблон
		tmpl, err := template.New("metrics").Parse(metricsTemplate)
		if err != nil {
			// В случае ошибки парсинга шаблона, отправляем 500 ошибку
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Устанавливаем заголовок для HTML контента
		c.Header("Content-Type", "text/html")

		// Статус 200 OK и отправка сгенерированного HTML-шаблона
		c.Writer.WriteHeader(http.StatusOK)
		err = tmpl.Execute(c.Writer, metrics)
		if err != nil {
			// В случае ошибки при рендеринге шаблона
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error rendering the template"})
		}
	}
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
		<li>{{.ID}}: {{.Value}}</li>
	{{else}}
		<li>No metrics available</li>
	{{end}}
	</ul>
</body>
</html>`
