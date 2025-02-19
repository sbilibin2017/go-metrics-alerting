package handlers

import (
	"go-metrics-alerting/internal/formatters"
	"go-metrics-alerting/internal/responders"

	"go-metrics-alerting/internal/templates"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/internal/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsService определяет контракт для обновления метрик.
type UpdateMetricsService interface {
	Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error)
}

// UpdateMetricsBodyHandler обрабатывает JSON-запросы на обновление метрик.
func UpdateMetricsBodyHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UpdateMetricsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}
		if err := validators.ValidateEmptyString(req.ID); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusNotFound, err.Error())
			return
		}
		if err := validators.ValidateMType(req.MType); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}
		if err := validators.ValidateDelta(req.MType, req.Delta); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}
		if err := validators.ValidateValue(req.MType, req.Value); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := service.Update(&req)
		if err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}

		responders.Respond(c, responders.JSONResponder, http.StatusOK, resp)
	}
}

// UpdateMetricsPathHandler обрабатывает обновление метрики через URL-параметры.
func UpdateMetricsPathHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mtype := c.Param("mtype")
		id := c.Param("id")
		valueStr := c.Param("value")

		var req types.UpdateMetricsRequest
		req.MType = types.MType(mtype)
		req.ID = id

		if err := validators.ValidateEmptyString(req.ID); err != nil {
			responders.Respond(c, responders.StringResponder, http.StatusNotFound, err.Error())
			return
		}

		if err := validators.ValidateMType(req.MType); err != nil {
			responders.Respond(c, responders.StringResponder, http.StatusBadRequest, err.Error())
			return
		}

		switch mtype {
		case string(types.Counter):
			if err := validators.ValidateCounterValue(valueStr); err != nil {
				responders.Respond(c, responders.StringResponder, http.StatusBadRequest, err.Error())
				return
			}
			value, _ := formatters.ParseInt64(valueStr)
			req.Delta = &value
		case string(types.Gauge):
			if err := validators.ValidateGaugeValue(valueStr); err != nil {
				responders.Respond(c, responders.StringResponder, http.StatusBadRequest, err.Error())
				return
			}
			value, _ := formatters.ParseFloat64(valueStr)
			req.Value = &value
		}

		_, err := service.Update(&req)
		if err != nil {
			responders.Respond(c, responders.StringResponder, http.StatusBadRequest, err.Error())
			return
		}

		responders.Respond(c, responders.StringResponder, http.StatusOK, "OK")
	}
}

// GetMetricValueService определяет контракт для получения значения метрики.
type GetMetricValueService interface {
	GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error)
}

// GetMetricValueBodyHandler обрабатывает JSON-запросы на получение метрик.
func GetMetricValueBodyHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.GetMetricValueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}
		if err := validators.ValidateEmptyString(req.ID); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusNotFound, err.Error())
			return
		}
		if err := validators.ValidateMType(req.MType); err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := service.GetMetricValue(&req)
		if err != nil {
			responders.Respond(c, responders.JSONResponder, http.StatusInternalServerError, err.Error())
			return
		}

		responders.Respond(c, responders.JSONResponder, http.StatusOK, resp)
	}
}

// GetMetricValuePathHandler обрабатывает получение метрик через URL-параметры.
func GetMetricValuePathHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mtype := c.Param("mtype")
		id := c.Param("id")

		if err := validators.ValidateMType(types.MType(mtype)); err != nil {
			responders.Respond(c, responders.StringResponder, http.StatusBadRequest, err.Error())
			return
		}

		req := &types.GetMetricValueRequest{
			MType: types.MType(mtype),
			ID:    id,
		}

		resp, err := service.GetMetricValue(req)
		if err != nil {
			responders.Respond(c, responders.StringResponder, http.StatusNotFound, err.Error())
			return
		}

		responders.Respond(c, responders.StringResponder, http.StatusOK, resp.Value)
	}
}

// GetAllMetricValuesService определяет контракт для получения всех метрик.
type GetAllMetricValuesService interface {
	GetAllMetricValues() []*types.GetMetricValueResponse
}

// GetAllMetricValuesHandler обрабатывает запрос и возвращает HTML-страницу со списком метрик.
func GetAllMetricValuesHandler(service GetAllMetricValuesService) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := service.GetAllMetricValues()
		c.Set("metrics", metrics)
		responders.Respond(c, responders.HTMLResponder, http.StatusOK, templates.MetricsTemplate)
	}
}
