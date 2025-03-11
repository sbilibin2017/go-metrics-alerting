package routers

import (
	"go-metrics-alerting/internal/apis/middlewares"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// MetricHandlers интерфейс для всех обработчиков метрик.
type MetricHandler interface {
	UpdateMetricWithPath(w http.ResponseWriter, r *http.Request)
	UpdateMetricWithBody(w http.ResponseWriter, r *http.Request)
	GetMetricWithPath(w http.ResponseWriter, r *http.Request)
	GetMetricWithBody(w http.ResponseWriter, r *http.Request)
	GetAllMetrics(w http.ResponseWriter, r *http.Request)
}

func NewMetricRouter(h MetricHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.TimeoutMiddleware(5 * time.Second))

	r.Post("/update/{type}/{id}/{value}", h.UpdateMetricWithPath)
	r.With(middlewares.GzipMiddleware).Post("/update/", h.UpdateMetricWithBody)

	r.Post("/value/{type}/{id}", h.GetMetricWithPath)
	r.With(middlewares.GzipMiddleware).Post("/value/", h.GetMetricWithBody)

	r.Get("/", h.GetAllMetrics)
	return r
}
