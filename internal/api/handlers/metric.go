package handlers

import (
	"net/http"

	"go-metrics-alerting/internal/api/handlers/responders"
	"go-metrics-alerting/internal/api/handlers/templates"
	"go-metrics-alerting/internal/api/types"
	"go-metrics-alerting/internal/api/validators"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsService определяет контракт для обновления метрик.
type UpdateMetricsService interface {
	Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error)
}

// UpdateMetricsHandler возвращает обработчик, принимающий сервис обновления метрик.
func UpdateMetricsHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UpdateMetricsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		if err := validators.ValidateID(req.ID); err != nil {
			responders.RespondWithError(c, http.StatusNotFound, err)
			return
		}
		if err := validators.ValidateMType(req.MType); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		if err := validators.ValidateDelta(req.MType, req.Delta); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		if err := validators.ValidateValue(req.MType, req.Value); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		resp, err := service.Update(&req)
		if err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		responders.RespondWithSuccess(c, http.StatusOK, resp)
	}
}

// GetMetricValueService определяет контракт для получения значения метрики.
type GetMetricValueService interface {
	GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error)
}

// GetMetricValueHandler возвращает обработчик, принимающий сервис получения метрик.
func GetMetricValueHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.GetMetricValueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		if err := validators.ValidateID(req.ID); err != nil {
			responders.RespondWithError(c, http.StatusNotFound, err)
			return
		}
		if err := validators.ValidateMType(req.MType); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		resp, err := service.GetMetricValue(&req)
		if err != nil {
			responders.RespondWithError(c, http.StatusInternalServerError, err)
			return
		}
		responders.RespondWithSuccess(c, http.StatusOK, resp)
	}
}

// GetAllMetricValuesService интерфейс для сервиса работы с метриками.
type GetAllMetricValuesService interface {
	GetAllMetricValues() []*types.GetMetricValueResponse
}

// GetAllMetricValuesHandler обрабатывает запрос и возвращает HTML-страницу со списком метрик.
func GetAllMetricValuesHandler(service GetAllMetricValuesService) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := service.GetAllMetricValues()
		responders.RespondWithHTML(c, http.StatusOK, templates.MetricsTemplate, metrics)
	}
}
