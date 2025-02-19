package middlewares

import (
	"go-metrics-alerting/internal/server/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	JSONContentType         string = "application/json"
	ErrUnsupportedMediaType string = "Unsupported Media Type"
)

// JSONContentTypeMiddleware проверяет, что Content-Type = application/json для POST-запросов и логирует информацию о запросах
func JSONContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			if c.ContentType() != JSONContentType {
				logger.Logger.Warn("Unsupported Media Type", zap.String("url", c.Request.URL.String()))
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": ErrUnsupportedMediaType})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
