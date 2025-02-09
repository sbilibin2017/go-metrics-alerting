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

func RegisterGetAllMetricValuesHandler(r *gin.Engine, svc GetAllValuesService) {
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
		{{end}}
		</ul>
	</body>
	</html>`

	// Парсим HTML-шаблон с обработкой ошибок
	tmpl, err := template.New("metrics").Parse(metricsTemplate)
	if err != nil {
		// В случае ошибки парсинга, возвращаем ошибку в лог
		panic("Error parsing template: " + err.Error())
	}

	// Регистрация шаблона для использования в рендере
	r.SetHTMLTemplate(tmpl) // Регистрация шаблона в gin

	r.GET("/", func(c *gin.Context) {
		getAllMetricValuesHandler(svc, c)
	})
}

// getAllMetricValuesHandler — обработчик получения всех метрик в HTML
func getAllMetricValuesHandler(service GetAllValuesService, c *gin.Context) {
	// Получаем все метрики через сервис
	metrics := service.GetAllMetricValues(c)

	// Проверка, что метрики не пустые
	if metrics == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No metrics found"})
		return
	}

	// Установка заголовков для HTML ответа
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))

	// Отправка HTML с метками
	c.HTML(http.StatusOK, "metrics", metrics)
}
