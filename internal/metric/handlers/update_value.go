package handlers

import (
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	emptyUpdateStringRequest string = ""
)

type UpdateValueValidator interface {
	Validate() error
}

type UpdateMetricValueRequest struct {
	Type  string
	Name  string
	Value string
}

func (r *UpdateMetricValueRequest) Validate() error {
	if r.Type == emptyUpdateStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric type is required",
		}
	}

	if r.Name == emptyUpdateStringRequest {
		return &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: "Metric name is required",
		}
	}

	if r.Value == emptyUpdateStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric value is required",
		}
	}

	return nil
}

type UpdateValueService interface {
	UpdateMetricValue(req *UpdateMetricValueRequest) error
}

func RegisterUpdateValueHandler(r *gin.Engine, svc UpdateValueService) {
	r.RedirectTrailingSlash = false

	r.POST("/update/:type/:name/:value", func(c *gin.Context) {
		updateValueHandler(svc, c)
	})
}

func updateValueHandler(service UpdateValueService, c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	updateRequest := &UpdateMetricValueRequest{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}

	// Проверяем валидацию запроса
	if err := updateRequest.Validate(); err != nil {
		apiErr, ok := err.(*apierror.APIError)
		if ok {
			// Используем c.String для ошибки валидации
			c.String(apiErr.Code, apiErr.Message)
		} else {
			// Если ошибка не APIError, возвращаем 500
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	// Вызываем сервис для обновления метрики
	if err := service.UpdateMetricValue(updateRequest); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Ответ об успешном обновлении
	response := "Metric updated"
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.Header("Content-Length", strconv.Itoa(len(response)))
	c.String(http.StatusOK, response)
}
