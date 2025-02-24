package routers

import (
	"go-metrics-alerting/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// RegisterMetricsHandlers регистрирует обработчики для работы с метриками.
func RegisterMetricsHandlers(
	r chi.Router,
	logger *zap.Logger,
	updateBodyHandler http.HandlerFunc,
	updatePathHandler http.HandlerFunc,
	getBodyHandler http.HandlerFunc,
	getPathHandler http.HandlerFunc,
	getAllHandler http.HandlerFunc,
) {
	// Регистрируем обработчики для каждого маршрута с нужными middlewares

	// POST для обновления метрики через тело запроса с JSON middleware
	r.With(
		middlewares.DateMiddleware,
		middlewares.ContentLengthMiddleware,
		middlewares.JSONContentType,
		middlewares.LoggingMiddleware(logger),
		middlewares.GzipMiddleware,
	).Post("/update/", updateBodyHandler)

	// POST для обновления метрики по ID с TextPlainContentType middleware
	r.With(
		middlewares.DateMiddleware,
		middlewares.ContentLengthMiddleware,
		middlewares.TextPlainContentType,
		middlewares.LoggingMiddleware(logger),
	).Post("/update/{type}/{id}/{value}", updatePathHandler)

	// GET для получения метрики через тело запроса с JSON middleware
	r.With(
		middlewares.DateMiddleware,
		middlewares.ContentLengthMiddleware,
		middlewares.JSONContentType,
		middlewares.LoggingMiddleware(logger),
		middlewares.GzipMiddleware,
	).Get("/value/", getBodyHandler)

	// GET для получения метрики по ID с TextPlainContentType middleware
	r.With(
		middlewares.DateMiddleware,
		middlewares.ContentLengthMiddleware,
		middlewares.TextPlainContentType,
		middlewares.LoggingMiddleware(logger),
	).Get("/value/{type}/{id}", getPathHandler)

	// GET для получения всех метрик с JSON middleware
	r.With(
		middlewares.DateMiddleware,
		middlewares.ContentLengthMiddleware,
		middlewares.HTMLContentType,
		middlewares.LoggingMiddleware(logger),
	).Get("/", getAllHandler)
}
