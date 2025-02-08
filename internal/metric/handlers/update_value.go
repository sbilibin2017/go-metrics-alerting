package handlers

import (
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const emptyUpdateStringRequest string = ""

// Валидатор входных данных
type UpdateValueValidator interface {
	Validate() error
}

type UpdateMetricValueRequest struct {
	Type  string
	Name  string
	Value string
}

func (r *UpdateMetricValueRequest) Validate() error {
	// Проверяем, указан ли тип метрики
	if r.Type == emptyUpdateStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric type is required",
		}
	}

	// Проверяем поддержку типа метрики
	if r.Type != "gauge" && r.Type != "counter" {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Unsupported metric type",
		}
	}

	// Проверяем, указано ли имя метрики
	if r.Name == emptyUpdateStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric name is required",
		}
	}

	// Проверяем корректность значения метрики
	if r.Value == emptyUpdateStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric value is required",
		}
	}

	if _, err := strconv.ParseFloat(r.Value, 64); err != nil {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid metric value",
		}
	}

	return nil
}

// Интерфейс сервиса обновления метрик
type UpdateValueService interface {
	UpdateMetricValue(req *UpdateMetricValueRequest) error
}

// Регистрация обработчиков обновления метрик
func RegisterUpdateValueHandler(r *gin.Engine, svc UpdateValueService) {
	r.RedirectTrailingSlash = false

	// Основной маршрут обновления метрик
	r.POST("/update/:type/:name/:value", func(c *gin.Context) {
		metricType := c.Param("type")
		metricName := c.Param("name")
		metricValue := c.Param("value")

		// Проверка на отсутствие важных параметров
		if metricType == "" || metricName == "" || metricValue == "" {
			c.String(http.StatusNotFound, "Required parameters missing (type, name, value)")
			return
		}

		updateValueHandler(svc, c)
	})

	// Обработка запросов с отсутствующим `name` или `value` и другие ошибки
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Route not found")
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

	// Валидация запроса
	if err := updateRequest.Validate(); err != nil {
		apiErr, ok := err.(*apierror.APIError)
		if ok {
			c.String(apiErr.Code, apiErr.Message)
		}
		return
	}

	// Обновляем значение метрики
	if err := service.UpdateMetricValue(updateRequest); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Успешный ответ
	response := "Metric updated"
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.Header("Content-Length", strconv.Itoa(len(response)))
	c.String(http.StatusOK, response)
}
