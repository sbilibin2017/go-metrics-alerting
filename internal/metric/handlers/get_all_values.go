package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricResponse struct {
	Type  string
	Name  string
	Value string
}

// Интерфейс для получения значения метрики
type GetAllValuesService interface {
	GetAllMetricValues() []*MetricResponse
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

	tmpl, _ := template.New("metrics").Parse(metricsTemplate)

	r.SetHTMLTemplate(tmpl) // Регистрация шаблона в gin

	r.GET("/", func(c *gin.Context) {
		getAllMetricValuesHandler(svc, c)
	})
}

// getAllMetricValuesHandler — обработчик получения всех метрик в HTML
func getAllMetricValuesHandler(service GetAllValuesService, c *gin.Context) {
	metrics := service.GetAllMetricValues()

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.HTML(http.StatusOK, "metrics", metrics)
}
