package routers

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MetricHandlers defines the set of handler methods required for managing metrics.
type MetricHandlers interface {
	UpdateMetricPathHandler(w http.ResponseWriter, r *http.Request)
	UpdatesMetricBodyHandler(w http.ResponseWriter, r *http.Request)
	UpdateMetricBodyHandler(w http.ResponseWriter, r *http.Request)
	GetMetricByTypeAndIDPathHandler(w http.ResponseWriter, r *http.Request)
	GetMetricByTypeAndIDBodyHandler(w http.ResponseWriter, r *http.Request)
	ListMetricsHTMLHandler(w http.ResponseWriter, r *http.Request)
}

type MetricRouter struct {
	*chi.Mux
	config *configs.ServerConfig
}

// NewMetricRouter initializes and returns a new MetricRouter with the provided handlers and config.
func NewMetricRouter(config *configs.ServerConfig, h MetricHandlers) *MetricRouter {
	r := chi.NewRouter()

	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.GzipMiddleware())
	r.Use(middlewares.TimeoutMiddleware())

	r.Post("/update/{type}/{id}/{value}", h.UpdateMetricPathHandler)
	r.Post("/updates/", h.UpdatesMetricBodyHandler)
	r.Post("/update/", h.UpdateMetricBodyHandler)
	r.Get("/value/{type}/{id}", h.GetMetricByTypeAndIDPathHandler)
	r.Post("/value/", h.GetMetricByTypeAndIDBodyHandler)
	r.Get("/", h.ListMetricsHTMLHandler)

	return &MetricRouter{Mux: r, config: config}
}
