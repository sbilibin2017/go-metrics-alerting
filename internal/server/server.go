package server

import (
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/middlewares"
	"go-metrics-alerting/internal/routers"

	"github.com/gin-gonic/gin"
)

// NewServerFactory создает новый сервер с роутером и middleware
func NewServer() *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = false

	r.Use(middlewares.JSONContentTypeMiddleware())
	r.Use(middlewares.LoggerMiddleware(logger.Logger))

	routers.RegisterMetricRoutes(
		r,
	)

	return r
}
