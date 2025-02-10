package handlers

import (
	"context"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Сервис обновления метрики
type UpdateValueService interface {
	UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error
}

// Валидационные интерфейсы
type UpdateValueMetricTypeValidator interface {
	Validate(metricType string) error
}

type UpdateValueMetricNameValidator interface {
	Validate(metricName string) error
}

type UpdateValueMetricValueValidator interface {
	Validate(metricValue string) error
}

type UpdateValueMetricGaugeValidator interface {
	Validate(metricValue string) error
}

type UpdateValueMetricCounterValidator interface {
	Validate(metricValue string) error
}

// Структура обработчика для обновления метрики
type UpdateValueHandler struct {
	service               UpdateValueService
	metricTypeValidator   UpdateValueMetricTypeValidator
	metricNameValidator   UpdateValueMetricNameValidator
	metricValueValidator  UpdateValueMetricValueValidator
	gaugeValueValidator   UpdateValueMetricGaugeValidator
	counterValueValidator UpdateValueMetricCounterValidator
}

// Регистрация маршрутов для обновления метрики
func RegisterUpdateMetricValueHandler(r *gin.Engine, service UpdateValueService,
	metricTypeValidator UpdateValueMetricTypeValidator,
	metricNameValidator UpdateValueMetricNameValidator,
	metricValueValidator UpdateValueMetricValueValidator,
	gaugeValueValidator UpdateValueMetricGaugeValidator,
	counterValueValidator UpdateValueMetricCounterValidator,
) {

	handler := &UpdateValueHandler{
		service:               service,
		metricTypeValidator:   metricTypeValidator,
		metricNameValidator:   metricNameValidator,
		metricValueValidator:  metricValueValidator,
		gaugeValueValidator:   gaugeValueValidator,
		counterValueValidator: counterValueValidator,
	}

	r.RedirectTrailingSlash = false

	// Маршрут для обновления метрики с типом, именем и значением
	r.POST("/update/:type/:name/:value", func(c *gin.Context) {
		handler.updateValueHandler(c)
	})

	// Обработчик на случай, если маршрут не найден
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Route not found")
	})
}

// Обработчик обновления метрики
func (h *UpdateValueHandler) updateValueHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	// Применяем валидатор типа метрики
	if err := h.metricTypeValidator.Validate(metricType); err != nil {
		// Вернуть ошибку 400, если тип метрики невалиден
		c.String(http.StatusBadRequest, err.Error()) // Ошибка типа метрики
		return
	}

	// Валидация имени метрики
	if err := h.metricNameValidator.Validate(metricName); err != nil {
		// Возвращаем 404, если ошибка валидации имени метрики
		c.String(http.StatusNotFound, err.Error())
		return
	}

	// Валидация значения метрики
	if err := h.metricValueValidator.Validate(metricValue); err != nil {
		c.String(http.StatusBadRequest, err.Error()) // Ошибка значения метрики
		return
	}

	// Дополнительная валидация значений в зависимости от типа метрики
	switch metricType {
	case string(types.Gauge):
		if err := h.gaugeValueValidator.Validate(metricValue); err != nil {
			c.String(http.StatusBadRequest, err.Error()) // Ошибка для Gauge
			return
		}
	case string(types.Counter):
		if err := h.counterValueValidator.Validate(metricValue); err != nil {
			c.String(http.StatusBadRequest, err.Error()) // Ошибка для Counter
			return
		}
	default:
		// Если тип метрики неизвестен, возвращаем ошибку 400 с описанием
		c.String(http.StatusBadRequest, errors.ErrUnsupportedMetricType.Error())
		return
	}

	// Подготовка запроса для обновления метрики
	updateRequest := &types.UpdateMetricValueRequest{
		Type:  types.MetricType(metricType),
		Name:  metricName,
		Value: metricValue,
	}

	// Обновляем метрику через сервис
	err := h.service.UpdateMetricValue(c, updateRequest)
	if err != nil {
		// Если ошибка типа *apierror.APIError, используем ее код и сообщение
		if apiErr, ok := err.(*apierror.APIError); ok {
			c.String(apiErr.Code, apiErr.Message) // Ошибка API
			return
		}
		// Внутренняя ошибка сервера
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Успешный ответ
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
	c.String(http.StatusOK, "Metric updated") // Gin сам установит Content-Length
}
