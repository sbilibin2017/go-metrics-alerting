package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// MetricServiceInterface defines methods for managing metrics.
type MetricService interface {
	UpdatesMetric(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
	GetMetricByTypeAndID(ctx context.Context, id types.MetricID) (*types.Metrics, error)
	ListAllMetrics(ctx context.Context) ([]*types.Metrics, error)
}

// MetricHandler contains the reference to the metric service.
type MetricHandler struct {
	svc MetricService
}

// NewMetricHandler creates a new instance of MetricHandler.
func NewMetricHandler(svc MetricService) *MetricHandler {
	return &MetricHandler{svc: svc}
}

// UpdateMetricPathHandler handles metric updates through path parameters.
func (h *MetricHandler) UpdateMetricPathHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricID := chi.URLParam(r, "id")
	metricValue := chi.URLParam(r, "value")

	// Log the start of request processing
	fmt.Printf("Handling update for metric: type=%s, id=%s, value=%s\n", metricType, metricID, metricValue)

	// Check if metric ID is provided
	if metricID == "" {
		http.Error(w, "Metric ID not found", http.StatusNotFound)
		fmt.Printf("Error: Metric ID not found: %s\n", metricID)
		return
	}

	// Convert value to the appropriate type
	var value *float64
	var delta *int64
	if metricType == string(types.Gauge) {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid value for gauge", http.StatusBadRequest)
			fmt.Printf("Error: Invalid value for gauge: %s, %v\n", metricValue, err)
			return
		}
		value = &val
	} else if metricType == string(types.Counter) {
		deltaVal, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid value for counter", http.StatusBadRequest)
			fmt.Printf("Error: Invalid value for counter: %s, %v\n", metricValue, err)
			return
		}
		delta = &deltaVal
	} else {
		http.Error(w, "Unknown metric type", http.StatusBadRequest)
		fmt.Printf("Error: Unknown metric type: %s\n", metricType)
		return
	}

	// Create the metric based on the inputs
	metric := &types.Metrics{
		Type:  metricType,
		ID:    metricID,
		Value: value,
		Delta: delta,
	}

	// Update the metric using the service
	updatedMetrics, err := h.svc.UpdatesMetric(r.Context(), []*types.Metrics{metric})
	if err != nil {
		http.Error(w, "Failed to update metric", http.StatusInternalServerError)
		fmt.Printf("Error: Failed to update metric: %s, %v\n", metricID, err)
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if len(updatedMetrics) > 0 {
		fmt.Fprintf(w, "Metric %s updated successfully", metric.ID)
	} else {
		fmt.Fprintln(w, "Failed to update the metric")
	}
}

// UpdatesMetricBodyHandler handles updating a list of metrics via request body.
func (h *MetricHandler) UpdatesMetricBodyHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []*types.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Error: Invalid input: %v\n", err)
		return
	}

	// Check that all metrics have valid IDs and values/deltas
	for _, metric := range metrics {
		if metric.ID == "" {
			http.Error(w, "Metric ID not found", http.StatusNotFound)
			fmt.Printf("Error: Metric ID not found: %s\n", metric.ID)
			return
		}

		if metric.Type == string(types.Gauge) {
			if metric.Value == nil {
				http.Error(w, "Missing value for gauge metric", http.StatusBadRequest)
				fmt.Printf("Error: Missing value for gauge metric: %s\n", metric.ID)
				return
			}
		} else if metric.Type == string(types.Counter) {
			if metric.Delta == nil {
				http.Error(w, "Missing delta for counter metric", http.StatusBadRequest)
				fmt.Printf("Error: Missing delta for counter metric: %s\n", metric.ID)
				return
			}
		} else {
			http.Error(w, "Unknown metric type", http.StatusBadRequest)
			fmt.Printf("Error: Unknown metric type: %s\n", metric.ID)
			return
		}
	}

	// Update the metrics using the service
	updatedMetrics, err := h.svc.UpdatesMetric(r.Context(), metrics)
	if err != nil {
		http.Error(w, "Failed to update metrics", http.StatusInternalServerError)
		fmt.Printf("Error: Failed to update metrics: %v\n", err)
		return
	}

	// Return the updated metrics in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMetrics)
}

// UpdateMetricBodyHandler handles updating a metric via request body.
func (h *MetricHandler) UpdateMetricBodyHandler(w http.ResponseWriter, r *http.Request) {
	var metric types.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Error: Invalid input: %v\n", err)
		return
	}

	// Check the metric ID
	if metric.ID == "" {
		http.Error(w, "Metric ID not found", http.StatusNotFound)
		fmt.Printf("Error: Metric ID not found: %s\n", metric.ID)
		return
	}

	// Validate the metric type and its value/delta
	if metric.Type == string(types.Gauge) {
		if metric.Value == nil {
			http.Error(w, "Missing value for gauge", http.StatusBadRequest)
			fmt.Printf("Error: Missing value for gauge: %s\n", metric.ID)
			return
		}
	} else if metric.Type == string(types.Counter) {
		if metric.Delta == nil {
			http.Error(w, "Missing delta for counter", http.StatusBadRequest)
			fmt.Printf("Error: Missing delta for counter: %s\n", metric.ID)
			return
		}
	} else {
		http.Error(w, "Unknown metric type", http.StatusBadRequest)
		fmt.Printf("Error: Unknown metric type: %s\n", metric.ID)
		return
	}

	// Update the metric using the service
	updatedMetrics, err := h.svc.UpdatesMetric(r.Context(), []*types.Metrics{&metric})
	if err != nil {
		http.Error(w, "Failed to update metric", http.StatusInternalServerError)
		fmt.Printf("Error: Failed to update metric: %s, %v\n", metric.ID, err)
		return
	}

	// Return the updated metric
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMetrics)
}

// GetMetricByTypeAndIDPathHandler handles retrieving a metric via path parameters.
func (h *MetricHandler) GetMetricByTypeAndIDPathHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricID := chi.URLParam(r, "id")

	// Log the request
	fmt.Printf("Fetching metric by type and ID: type=%s, id=%s\n", metricType, metricID)

	if metricID == "" {
		http.Error(w, "Metric ID not found", http.StatusNotFound)
		fmt.Printf("Error: Metric ID not found: %s\n", metricID)
		return
	}

	metric, err := h.svc.GetMetricByTypeAndID(r.Context(), types.MetricID{Type: metricType, ID: metricID})
	if err != nil {
		if errors.Is(err, services.ErrMetricNotFound) {
			http.Error(w, "Metric not found", http.StatusNotFound)
			fmt.Printf("Error: Metric not found: %s\n", metricID)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Error: Internal server error: %s, %v\n", metricID, err)
		return
	}

	// Prepare the value based on the metric type
	var value string
	if metric.Type == string(types.Gauge) && metric.Value != nil {
		value = fmt.Sprintf("%f", *metric.Value)
	} else if metric.Type == string(types.Counter) && metric.Delta != nil {
		value = fmt.Sprintf("%d", *metric.Delta)
	}

	// Return the metric
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))

}

// GetMetricByTypeAndIDBodyHandler handles retrieving a metric via request body.
func (h *MetricHandler) GetMetricByTypeAndIDBodyHandler(w http.ResponseWriter, r *http.Request) {
	var metricRequest types.MetricID
	if err := json.NewDecoder(r.Body).Decode(&metricRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Error: Invalid input: %v\n", err)
		return
	}

	// Log the request
	fmt.Printf("Fetching metric by body: %+v\n", metricRequest)

	if metricRequest.ID == "" {
		http.Error(w, "Metric ID not found", http.StatusNotFound)
		fmt.Printf("Error: Metric ID not found: %s\n", metricRequest.ID)
		return
	}

	metric, err := h.svc.GetMetricByTypeAndID(r.Context(), metricRequest)
	if err != nil {
		if errors.Is(err, services.ErrMetricNotFound) {
			http.Error(w, "Metric not found", http.StatusNotFound)
			fmt.Printf("Error: Metric not found: %s\n", metricRequest.ID)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Error: Internal server error: %s, %v\n", metricRequest.ID, err)
		return
	}

	// Return the metric as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metric)
}

// ListMetricsHTMLHandler returns the list of all metrics in HTML format.
func (h *MetricHandler) ListMetricsHTMLHandler(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.svc.ListAllMetrics(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve metrics", http.StatusInternalServerError)
		fmt.Printf("Error: Failed to retrieve metrics: %v\n", err)
		return
	}

	// Prepare the view model for HTML rendering
	type MetricViewModel struct {
		ID    string
		Value string
	}

	var viewModel []MetricViewModel
	for _, metric := range metrics {
		var value string
		if metric.Type == string(types.Gauge) && metric.Value != nil {
			value = fmt.Sprintf("%f", *metric.Value)
		} else if metric.Type == string(types.Counter) && metric.Delta != nil {
			value = fmt.Sprintf("%d", *metric.Delta)
		} else {
			value = "N/A"
		}
		viewModel = append(viewModel, MetricViewModel{
			ID:    metric.ID,
			Value: value,
		})
	}

	// HTML template for displaying metrics
	tmpl := `
	<html>
		<head><title>Metrics List</title></head>
		<body>
			<h1>Metrics List</h1>
			<ul>
				{{range .}}
					<li>{{.ID}}: {{.Value}}</li>
				{{else}}
					<li>No metrics found.</li>
				{{end}}
			</ul>
		</body>
	</html>
	`

	// Create a new template and execute it
	t, err := template.New("metrics").Parse(tmpl)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		fmt.Printf("Error: Failed to load template: %v\n", err)
		return
	}

	// Render the template with the data
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := t.Execute(w, viewModel); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		fmt.Printf("Error: Failed to render template: %v\n", err)
	}
}
