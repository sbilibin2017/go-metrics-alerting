package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-metrics-alerting/internal/apis/converters"
	"go-metrics-alerting/internal/apis/requests"
	"go-metrics-alerting/internal/apis/responses"
	"go-metrics-alerting/internal/apis/validators"
	"go-metrics-alerting/internal/types"
	"net/http"
	"text/template"
)

// UpdateMetricServiceInterface is an interface for updating a single metric
type UpdateMetricService interface {
	UpdateMetric(ctx context.Context, metric *types.Metrics) (*types.Metrics, error)
}

// UpdatesMetricServiceInterface is an interface for updating multiple metrics
type UpdatesMetricService interface {
	UpdateMetrics(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
}

// GetMetricServiceInterface is an interface for retrieving a single metric
type GetMetricService interface {
	GetMetric(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

// GetAllMetricsServiceInterface is an interface for retrieving all metrics
type GetAllMetricsService interface {
	GetAllMetrics(ctx context.Context) []*types.Metrics
}

type MetricService interface {
	UpdateMetricService
	UpdatesMetricService
	GetMetricService
	GetAllMetricsService
}

type MetricHandler struct {
	*UpdateMetricWithPathHandler
	*UpdateMetricWithBodyHandler
	*GetMetricWithPathHandler
	*GetMetricWithBodyHandler
	*GetAllMetricsHandler
}

func NewMetricHandler(svc MetricService) *MetricHandler {
	return &MetricHandler{
		UpdateMetricWithPathHandler: &UpdateMetricWithPathHandler{svc},
		UpdateMetricWithBodyHandler: &UpdateMetricWithBodyHandler{svc},
		GetMetricWithPathHandler:    &GetMetricWithPathHandler{svc},
		GetMetricWithBodyHandler:    &GetMetricWithBodyHandler{svc},
		GetAllMetricsHandler:        &GetAllMetricsHandler{svc},
	}
}

// UpdateMetricWithPathHandler Implementation
type UpdateMetricWithPathHandler struct {
	service UpdateMetricService
}

func (h *UpdateMetricWithPathHandler) UpdateMetricWithPath(w http.ResponseWriter, r *http.Request) {
	id := requests.GetPathParam(r, "id")
	mtype := requests.GetPathParam(r, "type")
	value := requests.GetPathParam(r, "value")

	var metric types.Metrics

	err := validators.ValidateMetricID(id)
	if err != nil {
		responses.NotFoundResponse(w, err)
		return
	}
	metric.ID = id

	err = validators.ValidateMetricType(mtype)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}
	metric.Type = types.MetricType(mtype)

	switch metric.Type {
	case types.Counter:
		err = validators.ValidateCounterMetricString(value)
		if err != nil {
			responses.BadRequestResponse(w, err)
			return
		}
		d, _ := converters.ConvertStringToInt64(value)
		metric.Delta = &d
	case types.Gauge:
		err = validators.ValidateGaugeMetricString(value)
		if err != nil {
			responses.BadRequestResponse(w, err)
			return
		}
		v, _ := converters.ConvertStringToFloat64(value)
		metric.Value = &v
	}

	_, err = h.service.UpdateMetric(r.Context(), &metric)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	responses.TextResponse(w, "Metric updated successfully")
}

// UpdateMetricWithBodyHandler Implementation
type UpdateMetricWithBodyHandler struct {
	service UpdateMetricService
}

func (h *UpdateMetricWithBodyHandler) UpdateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	var metric types.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid json data", http.StatusInternalServerError)
		return
	}

	err := validators.ValidateMetricID(string(metric.ID))
	if err != nil {
		responses.NotFoundResponse(w, err)
		return
	}

	switch metric.Type {
	case types.Counter:
		err = validators.ValidateCounterMetric(metric.Delta)
		if err != nil {
			responses.BadRequestResponse(w, err)
			return
		}
	case types.Gauge:
		err = validators.ValidateGaugeMetric(metric.Value)
		if err != nil {
			responses.BadRequestResponse(w, err)
			return
		}
	}

	updatedMetric, err := h.service.UpdateMetric(r.Context(), &metric)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	responses.JsonResponse(w, updatedMetric)

}

// GetMetricWithPathHandler Implementation
type GetMetricWithPathHandler struct {
	service GetMetricService
}

func (h *GetMetricWithPathHandler) GetMetricWithPath(w http.ResponseWriter, r *http.Request) {
	// Extract type and id from the URL path
	mtype := requests.GetPathParam(r, "type")
	id := requests.GetPathParam(r, "id")

	// Validate the metric ID
	err := validators.ValidateMetricID(id)
	if err != nil {
		responses.NotFoundResponse(w, err)
		return
	}

	// Validate the metric type
	err = validators.ValidateMetricType(mtype)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	// Create a MetricID struct to pass to the service
	metricID := types.MetricID{
		ID:   id,
		Type: types.MetricType(mtype),
	}

	// Fetch the metric by ID and type
	metric, err := h.service.GetMetric(r.Context(), metricID)
	if err != nil {
		responses.NotFoundResponse(w, err)
		return
	}

	var value string
	switch metricID.Type {
	case types.Counter:
		value = fmt.Sprint(*metric.Delta)
	case types.Gauge:
		value = fmt.Sprint(*metric.Value)
	}
	responses.TextResponse(w, value)
}

// GetMetricWithBodyHandler Implementation
type GetMetricWithBodyHandler struct {
	service GetMetricService
}

func (h *GetMetricWithBodyHandler) GetMetricWithBody(w http.ResponseWriter, r *http.Request) {
	// Extract metric details from request body
	var metricRequest types.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricRequest); err != nil {
		http.Error(w, "invalid json data", http.StatusInternalServerError)
		return
	}

	// Validate the metric ID in the body
	err := validators.ValidateMetricID(string(metricRequest.ID))
	if err != nil {
		responses.NotFoundResponse(w, err)
		return
	}

	// Create a MetricID struct to pass to the service
	metricID := types.MetricID{
		ID:   string(metricRequest.ID),
		Type: metricRequest.Type,
	}

	// Fetch the metric by ID and type
	metric, err := h.service.GetMetric(r.Context(), metricID)
	if err != nil {
		responses.NotFoundResponse(w, err)
		return
	}

	// Return the metric as JSON
	responses.JsonResponse(w, metric)
}

// GetAllMetricsHandler Implementation
type GetAllMetricsHandler struct {
	service GetAllMetricsService
}

func (h *GetAllMetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	// Get all metrics from the service
	metrics := h.service.GetAllMetrics(r.Context())

	// Parse the HTML template
	tmpl, err := template.New("metrics").Parse(metricsTemplate)
	if err != nil {
		responses.InternalServerErrorResponse(w, err)
		return
	}

	// If no metrics are available, pass an empty slice to the template
	if metrics == nil || len(metrics) == 0 {
		metrics = []*types.Metrics{}
	}

	// Prepare a slice to store the formatted metric values
	var values []map[string]string
	for _, metric := range metrics {
		// Prepare the value to display based on the metric type
		var value string
		switch metric.Type {
		case types.Counter:
			if metric.Delta != nil {
				value = fmt.Sprintf("%d", *metric.Delta)
			} else {
				value = "0" // Handle if Delta is nil
			}
		case types.Gauge:
			if metric.Value != nil {
				value = fmt.Sprintf("%f", *metric.Value)
			} else {
				value = "0.0" // Handle if Value is nil
			}
		default:
			value = "Unknown" // In case of unsupported metric type
		}

		// Append the metric ID and its value to the values slice
		values = append(values, map[string]string{
			"ID":    metric.ID,
			"Value": value,
		})
	}

	// Render the template with the metrics data
	err = tmpl.Execute(w, values)
	if err != nil {
		responses.InternalServerErrorResponse(w, err)
		return
	}
}

// Константа для HTML-шаблона.
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
