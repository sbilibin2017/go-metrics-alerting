package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Интерфейс для получения значения метрики
type GetAllValuesService interface {
	GetAllMetricValues(ctx context.Context) []*types.MetricResponse
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

// Функция для регистрации обработчиков
func RegisterGetAllMetricValuesHandler(r *gin.Engine, svc GetAllValuesService) {
	// Загружаем шаблон
	tmpl, _ := template.New("metrics").Parse(metricsTemplate)

	// Устанавливаем шаблон для использования в Gin
	r.SetHTMLTemplate(tmpl)

	// Регистрируем обработчик на главную страницу
	r.GET("/", func(c *gin.Context) {
		renderMetricsPage(svc, c)
	})
}

// Обработчик рендеринга страницы с метками
func renderMetricsPage(service GetAllValuesService, c *gin.Context) {
	// Получаем метрики через сервис
	metrics := service.GetAllMetricValues(c)

	// Если метрики отсутствуют, передаем пустой срез
	if metrics == nil {
		metrics = []*types.MetricResponse{}
	}

	// Установка заголовков для HTML ответа
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))

	// Отправка HTML с метками (или с сообщением, что нет метрик)
	c.HTML(http.StatusOK, "metrics", metrics)
}
