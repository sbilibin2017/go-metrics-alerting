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
	// Если тип пустой, возвращаем ошибку с кодом 400
	if r.Type == emptyGetStringRequest {
		return &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Metric type is required",
		}
	}

	// Если имя пустое, возвращаем ошибку с кодом 400
	if r.Name == emptyGetStringRequest {
		return &apierror.APIError{
			Code:    http.StatusNotFound,
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
	// Отключаем автоматический редирект на маршрут с "/"
	r.RedirectTrailingSlash = false

	// Основной обработчик получения значения метрики
	r.GET("/value/:type/:name", func(c *gin.Context) {
		getMetricValueHandler(svc, c)
	})

	// Если передан только тип метрики без имени — возвращаем ошибку 400 с нужным текстом
	r.GET("/value/:type", func(c *gin.Context) {
		metricType := c.Param("type")

		// Создаем запрос
		getRequest := &GetMetricValueRequest{
			Type: metricType,
			Name: "", // Имя отсутствует
		}

		// Проверка валидации запроса
		if err := getRequest.Validate(); err != nil {
			// Если ошибка валидации, отправляем ошибку с кодом и сообщением
			apiErr, ok := err.(*apierror.APIError)
			if ok {
				c.String(apiErr.Code, apiErr.Message)
			}
			return
		}
	})
}

// Обработчик получения значения метрики
func getMetricValueHandler(service GetValueService, c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	// Создаем запрос
	getRequest := &GetMetricValueRequest{
		Type: metricType,
		Name: metricName,
	}

	// Проверка валидации запроса
	if err := getRequest.Validate(); err != nil {
		// Если ошибка валидации, отправляем ошибку с кодом и сообщением
		apiErr, ok := err.(*apierror.APIError)
		if ok {
			c.String(apiErr.Code, apiErr.Message)
		}
		return
	}

	// Получаем значение метрики через сервис
	metricValue, err := service.GetMetricValue(getRequest)
	if err != nil {
		// Если ошибка при получении метрики, отправляем соответствующий ответ
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message)
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	// Если метрика найдена, отправляем значение метрики
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.Header("Content-Length", strconv.Itoa(len(metricValue)))
	c.String(http.StatusOK, metricValue)
}
