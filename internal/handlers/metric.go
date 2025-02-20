package handlers

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/services"
	"go-metrics-alerting/internal/types"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateMetricsServiceInterface интерфейс для сервиса обновления метрик
type UpdateMetricsService interface {
	UpdateMetricValue(metric *domain.Metric) (*domain.Metric, error)
}

// GetMetricValueServiceInterface интерфейс для сервиса получения метрики по ID
type GetMetricValueService interface {
	GetMetricValue(id string, mType domain.MType) (*domain.Metric, error)
}

// GetAllMetricValuesServiceInterface интерфейс для сервиса получения всех метрик
type GetAllMetricValuesService interface {
	GetAllMetricValues() []*domain.Metric
}

// UpdateMetricsBodyHandler обрабатывает обновление метрики через body запроса
func UpdateMetricsBodyHandler(updateMetricService UpdateMetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.UpdateMetricBodyRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			handleError(c, err)
			return
		}
		if err := request.Validate(); err != nil {
			handleError(c, err)
			return
		}
		updatedMetric, err := updateMetricService.UpdateMetricValue(request.ToMetric())
		if err != nil {
			handleError(c, err)
			return
		}
		sendResponse(c, http.StatusOK, updatedMetric, "application/json", nil) // Для JSON-шаблон не нужен
	}
}

// // UpdateMetricsPathHandler обрабатывает обновление метрики через параметры в пути
// func UpdateMetricsPathHandler(updateMetricService UpdateMetricsService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var request types.UpdateMetricPathRequest
// 		request.Name = c.Param("name")
// 		request.Type = c.Param("type")
// 		request.Value = c.Param("value")
// 		if err := request.Validate(); err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		_, err := updateMetricService.UpdateMetricValue(request.ToMetric())
// 		if err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		sendResponse(c, http.StatusOK, "Metric is updated", "text/plain", nil) // Для plain текста тоже шаблон не нужен
// 	}
// }

// // GetMetricValueBodyHandler обрабатывает получение метрики через body запроса
// func GetMetricValueBodyHandler(getMetricValueService GetMetricValueService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var request types.GetMetricBodyRequest
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		if err := request.Validate(); err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		metric, err := getMetricValueService.GetMetricValue(request.ID, domain.MType(request.MType))
// 		if err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		sendResponse(c, http.StatusOK, &types.GetMetricBodyResponse{
// 			ID:    metric.ID,
// 			MType: string(metric.MType),
// 			Value: metric.Value,
// 		}, "application/json", nil) // Для JSON-шаблон не нужен
// 	}
// }

// // GetMetricValuePathHandler обрабатывает получение метрики через параметры в пути
// func GetMetricValuePathHandler(getMetricValueService GetMetricValueService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var request types.GetMetricBodyRequest

// 		request.ID = c.Param("id")
// 		request.MType = c.Param("mtype")
// 		if err := request.Validate(); err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		metric, err := getMetricValueService.GetMetricValue(request.ID, domain.MType(request.MType))
// 		if err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		sendResponse(c, http.StatusOK, metric.Value, "text/plain", nil) // Для plain текста тоже шаблон не нужен
// 	}
// }

// // GetAllMetricValuesHandler обрабатывает получение всех метрик и возвращает HTML страницу
// func GetAllMetricValuesHandler(getAllMetricValuesService GetAllMetricValuesService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		metrics := getAllMetricValuesService.GetAllMetricValues()

// 		var response []types.GetMetricBodyResponse
// 		for _, metric := range metrics {
// 			response = append(response, types.GetMetricBodyResponse{
// 				ID:    metric.ID,
// 				MType: string(metric.MType),
// 				Value: metric.Value,
// 			})
// 		}
// 		tmpl, err := template.New("metrics").Parse(templates.MetricsTemplate)
// 		if err != nil {
// 			handleError(c, err)
// 			return
// 		}
// 		sendResponse(c, http.StatusOK, response, "text/html", tmpl) // Передаем шаблон для HTML
// 	}
// }

// SetHeaders устанавливает заголовки для ответа.
func setHeaders(c *gin.Context, contentType string) {
	c.Header("Content-Type", contentType+"; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
}

// handleError обрабатывает ошибки и возвращает соответствующий ответ
func handleError(c *gin.Context, err error) {
	if err != nil {
		if err == services.ErrMetricNotFound {
			contentType := c.GetHeader("Accept")
			if contentType == "text/plain" {
				c.String(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
			}
		} else {
			contentType := c.GetHeader("Accept")
			if contentType == "text/plain" {
				c.String(http.StatusBadRequest, err.Error())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
}

// sendResponse отправляет ответ в зависимости от типа контента.
func sendResponse(c *gin.Context, statusCode int, body interface{}, contentType string, tmpl *template.Template) {
	setHeaders(c, contentType)

	// Отправляем ответ в формате JSON или текстовом, в зависимости от типа контента.
	if contentType == "application/json" {
		c.JSON(statusCode, body)
	} else if contentType == "text/plain" {
		c.String(statusCode, body.(string))
	} else if contentType == "text/html" {
		// Рендерим шаблон, если передан тип "text/html"
		if tmpl != nil {
			err := tmpl.Execute(c.Writer, body)
			if err != nil {
				handleError(c, err)
			}
		}
	}
}
