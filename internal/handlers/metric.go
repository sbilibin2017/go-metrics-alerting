package handlers

import (
	"go-metrics-alerting/internal/responders"
	"go-metrics-alerting/internal/templates"
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsServiceInterface интерфейс для сервиса обновления метрик
type UpdateMetricsService interface {
	UpdateMetricValue(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, *types.APIErrorResponse)
	ParseMetricValues(mtype, valueStr string) (*float64, *int64, *types.APIErrorResponse)
}

// GetMetricValueServiceInterface интерфейс для сервиса получения метрики по ID
type GetMetricValueService interface {
	GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, *types.APIErrorResponse)
}

// GetAllMetricValuesServiceInterface интерфейс для сервиса получения всех метрик
type GetAllMetricValuesService interface {
	GetAllMetricValues() ([]*types.GetMetricValueResponse, *types.APIErrorResponse)
}

// Обработчик для обновления метрик с параметрами пути (отправляем text/plain)
func UpdateMetricsPathHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mtype := c.Param("mtype")
		id := c.Param("id")
		valueStr := c.Param("value")

		value, delta, errResponse := service.ParseMetricValues(mtype, valueStr)
		if errResponse != nil {
			responders.SendErrorJSON(c, errResponse.Status, errResponse.Message)
			return
		}

		req := &types.UpdateMetricsRequest{
			ID:    id,
			MType: types.MType(mtype),
			Delta: delta,
			Value: value,
		}

		_, errResponse = service.UpdateMetricValue(req)
		if errResponse != nil {
			responders.SendErrorJSON(c, errResponse.Status, errResponse.Message)
			return
		}

		responders.SendSuccessText(c, http.StatusOK, "OK")
	}
}

// Обработчик для обновления метрик с телом запроса (отправляем JSON)
func UpdateMetricsBodyHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UpdateMetricsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.SendErrorJSON(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		response, errResponse := service.UpdateMetricValue(&req)
		if errResponse != nil {
			responders.SendErrorJSON(c, errResponse.Status, errResponse.Message)
			return
		}

		responders.SendSuccessJSON(c, http.StatusOK, response)
	}
}

// Обработчик для получения метрики по ID с телом запроса (отправляем JSON)
func GetMetricValueBodyHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.GetMetricValueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.SendErrorJSON(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		response, errResponse := service.GetMetricValue(&req)
		if errResponse != nil {
			responders.SendErrorJSON(c, errResponse.Status, errResponse.Message)
			return
		}

		responders.SendSuccessJSON(c, http.StatusOK, response)
	}
}

// Обработчик для получения метрики по ID с параметрами пути (отправляем text/plain)
func GetMetricValuePathHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mtype := c.Param("mtype")
		id := c.Param("id")

		req := &types.GetMetricValueRequest{
			ID:    id,
			MType: types.MType(mtype),
		}

		_, errResponse := service.GetMetricValue(req)
		if errResponse != nil {
			responders.SendErrorJSON(c, errResponse.Status, errResponse.Message)
			return
		}

		responders.SendSuccessText(c, http.StatusOK, "OK")
	}
}

// Обработчик для получения всех метрик (отправляем HTML)
func GetAllMetricValuesHandler(service GetAllMetricValuesService) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics, _ := service.GetAllMetricValues()
		responders.SendSuccessHTML(c, http.StatusOK, templates.MetricsTemplate, metrics)
	}
}
