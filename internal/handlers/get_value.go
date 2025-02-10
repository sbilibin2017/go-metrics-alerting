package handlers

import (
	"context"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateValueService интерфейс для сервисов обновления метрик.
type GetValueMetricTypeEmptyStringValidator interface {
	Validate(metricType string) error
}

type GetValueMetricNameEmptyStringValidator interface {
	Validate(metricName string) error
}

// Интерфейс для получения значения метрики
type GetValueService interface {
	GetMetricValue(ctx context.Context, req *types.GetMetricValueRequest) (string, error)
}

// UpdateValueHandler структура обработчика с внедрением зависимостей
type GetValueHandler struct {
	service             GetValueService
	metricTypeValidator GetValueMetricTypeEmptyStringValidator
	metricNameValidator GetValueMetricNameEmptyStringValidator
}

// Регистрация обработчика для получения значения метрики
func RegisterGetMetricValueHandler(r *gin.Engine, handler *GetValueHandler) {
	r.RedirectTrailingSlash = false

	r.GET("/value/:type/:name", func(c *gin.Context) {
		handler.getMetricValueHandler(c)
	})

	r.GET("/value/:type", func(c *gin.Context) {
		handler.getMetricValueByTypeHandler(c)
	})
}

// Обработчик получения значения метрики по типу и имени
func (h *GetValueHandler) getMetricValueHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	// Применяем валидаторы
	if err := h.metricTypeValidator.Validate(metricType); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.metricNameValidator.Validate(metricName); err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	getRequest := &types.GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}

	// Передаем контекст в сервис
	metricValueResp, err := h.service.GetMetricValue(c, getRequest)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	// Убеждаемся, что Content-Length устанавливается корректно
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.String(http.StatusOK, metricValueResp) // Gin сам установит Content-Length
}

// Обработчик получения значения метрики только по типу
func (h *GetValueHandler) getMetricValueByTypeHandler(c *gin.Context) {
	metricType := c.Param("type")

	// Применяем валидаторы
	if err := h.metricTypeValidator.Validate(metricType); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	getRequest := &types.GetMetricValueRequest{
		Type: metricType,
	}

	// Передаем контекст в сервис
	metricValue, err := h.service.GetMetricValue(c, getRequest) // доступ через h.service
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
		}
		return
	}

	c.String(http.StatusOK, metricValue) // Gin сам установит Content-Length
}
