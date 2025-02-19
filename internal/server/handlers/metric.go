package handlers

import (
	"go-metrics-alerting/internal/server/formatters"
	"go-metrics-alerting/internal/server/responders"
	"go-metrics-alerting/internal/server/templates"
	"go-metrics-alerting/internal/server/types"
	"go-metrics-alerting/internal/server/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsService определяет контракт для обновления метрик.
type UpdateMetricsService interface {
	Update(req *types.UpdateMetricsRequest) (*types.UpdateMetricsResponse, error)
}

// UpdateMetricsHandler возвращает обработчик, принимающий сервис обновления метрик.
func UpdateMetricsBodyHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UpdateMetricsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		if err := validators.ValidateEmptyString(req.ID); err != nil {
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

// UpdateMetricsPathHandler обновляет метрику с использованием данных, переданных через путь URL.
func UpdateMetricsPathHandler(service UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mtype := c.Param("mtype")
		id := c.Param("id")
		valueStr := c.Param("value")

		var req types.UpdateMetricsRequest
		req.MType = types.MType(mtype)
		req.ID = id

		// Валидация ID с использованием валидатора ValidateEmptyString
		if err := validators.ValidateEmptyString(req.ID); err != nil {
			responders.RespondWithError(c, http.StatusNotFound, err)
			return
		}

		// Валидация типа метрики
		if err := validators.ValidateMType(req.MType); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}

		// Обрабатываем значение в зависимости от типа метрики
		switch mtype {
		case string(types.Counter):
			if err := validators.ValidateCounterValue(valueStr); err != nil {
				responders.RespondWithError(c, http.StatusBadRequest, err)
				return
			}
			value, _ := formatters.ParseInt64(valueStr)
			req.Delta = &value
		case string(types.Gauge):
			if err := validators.ValidateGaugeValue(valueStr); err != nil {
				responders.RespondWithError(c, http.StatusBadRequest, err)
				return
			}
			value, _ := formatters.ParseFloat64(valueStr)
			req.Value = &value
		}

		// Вызываем сервис обновления метрики
		resp, err := service.Update(&req)
		if err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}

		// Отправляем успешный ответ
		responders.RespondWithSuccess(c, http.StatusOK, resp)
	}
}

// GetMetricValueService определяет контракт для получения значения метрики.
type GetMetricValueService interface {
	GetMetricValue(req *types.GetMetricValueRequest) (*types.GetMetricValueResponse, error)
}

// GetMetricValueHandler возвращает обработчик, принимающий сервис получения метрик.
func GetMetricValueBodyHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.GetMetricValueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		if err := validators.ValidateEmptyString(req.ID); err != nil {
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

func GetMetricValuePathHandler(service GetMetricValueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mtype := c.Param("mtype")
		id := c.Param("id")

		// Валидация типа метрики
		if err := validators.ValidateMType(types.MType(mtype)); err != nil {
			responders.RespondWithError(c, http.StatusBadRequest, err)
			return
		}

		req := &types.GetMetricValueRequest{
			MType: types.MType(mtype),
			ID:    id,
		}

		// Получение метрики
		resp, err := service.GetMetricValue(req)
		if err != nil {
			responders.RespondWithError(c, http.StatusNotFound, err)
			return
		}

		// Ответ с успешным статусом
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
