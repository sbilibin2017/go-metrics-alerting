package handlers

import (
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/internal/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateMetricService интерфейс для сервиса обновления метрики.
type UpdateMetricService interface {
	UpdateMetric(req *types.MetricsRequest) (*types.MetricsResponse, error)
}

// UpdateMetricHandler обновляет метрику, передавая интерфейс UpdateMetricService.
func UpdateMetricHandler(svc UpdateMetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Определяем метрику как указатель
		req := &types.MetricsRequest{}

		// Пробуем привязать JSON в тело запроса к структуре metric
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Валидация ID метрики
		if err := validators.ValidateID(req.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация типа метрики (MType)
		if err := validators.ValidateMType(types.MType(req.MType)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация Delta для типа Counter
		if err := validators.ValidateDelta(types.MType(req.MType), req.Delta); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация Value для типа Gauge
		if err := validators.ValidateValue(types.MType(req.MType), req.Value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Логика обновления метрики через сервис
		metric, err := svc.UpdateMetric(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update metric"})
			return
		}

		// Возвращаем обновленную метрику
		c.JSON(http.StatusOK, gin.H{"data": metric})
	}
}
