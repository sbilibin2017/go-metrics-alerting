package handlers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/templates"
	"go-metrics-alerting/internal/types"
	"html/template"
	"net/http"
)

// UpdateMetricService интерфейс для обновления метрики
type GetAllMetricsService interface {
	GetAllMetrics() []*domain.Metrics
}

// GetAllMetricsHandler - обработчик для получения всех метрик и отображения их на HTML-странице.
func GetAllMetricsHandler(service GetAllMetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("metrics").Parse(templates.MetricsTemplate)

		metrics := service.GetAllMetrics()

		var responseMetrics []types.GetAllMetricsResponse
		for _, metric := range metrics {
			resp := types.GetAllMetricsResponse{}
			responseMetrics = append(responseMetrics, resp.FromDomain(metric))
		}

		tmpl.Execute(w, responseMetrics)
	}
}
