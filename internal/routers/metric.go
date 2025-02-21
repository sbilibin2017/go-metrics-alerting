package routers

import (
	"github.com/gin-gonic/gin"
)

// RegisterMetricRoutes регистрирует маршруты для метрик
func RegisterMetricRoutes(
	router *gin.Engine,
	updateMetricsBodyHandler gin.HandlerFunc,
	updateMetricsPathHandler gin.HandlerFunc,
	getMetricValueBodyHandler gin.HandlerFunc,
	getMetricValuePathHandler gin.HandlerFunc,
	getAllMetricValuesHandler gin.HandlerFunc,
) {
	router.POST("/update/", updateMetricsBodyHandler)
	router.POST("/update/:mtype/:id/:value", updateMetricsPathHandler)
	router.POST("/value/", getMetricValueBodyHandler)
	router.GET("/value/:mtype/:id", getMetricValuePathHandler)
	router.GET("/", getAllMetricValuesHandler)
}
