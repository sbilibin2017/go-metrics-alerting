package handlers

import (
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Структура для запроса получения значения метрики
type GetMetricValueRequest struct {
	Type string
	Name string
}

const (
	emptyGetStringRequest string = ""
)

// Метод валидации для запроса
func (r *GetMetricValueRequest) Validate() error {
	if r.Type == emptyGetStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric type is required",
		}
	}

	if r.Name == emptyGetStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric name is required",
		}
	}

	return nil
}

// Интерфейс для получения значения метрики
type GetValueService interface {
	GetMetricValue(req *GetMetricValueRequest) (string, error)
}

// Регистрация обработчика для получения значения метрики
func RegisterGetMetricValueHandler(r *gin.Engine, svc GetValueService) {
	r.RedirectTrailingSlash = false

	r.GET("/value/:type/:name", func(c *gin.Context) {
		getMetricValueHandler(svc, c)
	})

	r.GET("/value/:type", getMetricValueByTypeHandler)
}

// Обработчик получения значения метрики
func getMetricValueHandler(service GetValueService, c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	getRequest := &GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}

	if err := getRequest.Validate(); err != nil {
		apiErr, ok := err.(*apierror.APIError)
		if ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusBadRequest, "Invalid metric request")
		}
		return
	}

	metricValue, err := service.GetMetricValue(getRequest)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.Header("Content-Length", strconv.Itoa(len(metricValue)))
	c.String(http.StatusOK, metricValue)
}

// Обработчик для получения значения метрики только по типу
func getMetricValueByTypeHandler(c *gin.Context) {
	metricType := c.Param("type")

	getRequest := &GetMetricValueRequest{
		Type: metricType,
		Name: "",
	}

	if err := getRequest.Validate(); err != nil {
		apiErr, ok := err.(*apierror.APIError)
		if ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusBadRequest, "Invalid metric request")
		}
		return
	}

	c.String(http.StatusBadRequest, "Metric name is required")
}
